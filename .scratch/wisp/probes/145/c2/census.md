# 145-c2 快照字段普查（只读腿 · 零 Go 命令）

> 派单：编排者直令（无独立 dispatch 文件，代号 `145-c2`）。工单＝`.scratch/wisp/issues/145-panel-snapshot-has-four-fields-so-eleven-of-the-fourteen-ui-states-have-no-input-to-render-grow-the-go-side-carrier.md`。
> 本程只读源码＋grep/find＋git log/show/diff，**一枚 go 命令未跑**；写面只到 `.scratch/wisp/probes/145/c2/**`；禁面一字未动。

## §0 起手锚

- 取数命令（同一发 shell）：`date "+%Y-%m-%d %H:%M:%S%z"` ＋ `git log -1 --format="%h %ci"` ＋ `git status --porcelain -- cmd internal tools docs scripts .scratch | wc -l`
- 读数（2026-10-03）：
  - `date` = `2026-10-03 09:26:36+0800`
  - HEAD = `5f9ff9d4 2026-10-03 09:25:59 +0800`
  - 脏项计数 = `365`
- 相关目录脏项（本轮 `git status --porcelain -- internal/panel/ cmd/wisp/ internal/config/ internal/agent/approval/ internal/tools/ internal/risk/` 现量）：
  - ` M cmd/wisp/config_reload.go`
  - ` M internal/tools/paths.go`
  - `?? cmd/wisp/config_readers_255.go`
  - `?? internal/config/tiers_app_255r2_test.go`
  - ⇒ **`internal/panel/**`（composer.go/pump.go/approval.go/子件）与 `cmd/wisp/run.go`、`cmd/wisp/panel_pump.go` 此刻均不在脏项里** ⇒ 本程读到的这些文件＝HEAD `5f9ff9d4` 已提交态，非他人未提交中间态。⚠ 在飞的三枚写腿（config/approval/tools+risk）可能随时改脏，复量者请以自己现量为准。

## §1 快照字段名册（三种"没做"分开标）

**三态口径**（本票票面定义）：①＝字段写了但生产里恒空＝没接；②＝真源存在但没被搬进快照＝差搬运；③＝真源压根不存在＝一块没写。
**⚠ 本票票面标题那句"恰四字段"已过期**：`Snapshot` 现量 **6 枚直接字段**（`internal/panel/composer.go:57-92`，非票面/旧普查的 `:43-49`／`:44-49`），`ComposerState` 现量 **11 枚字段**（`composer.go:235-270`，非票面 AC#2b 的"6 枚"、非 r2 的"7→7"）。增长来自后续票：181（git）、145-r3（currentModel/modelKnown）、248（credential）、197（tasks）、200（instructions）。

**生产构造链已存在（推翻票面 ⑤"从来没有生产代码构造过"那一支）**：唯一装配根＝`cmd/wisp/run.go:699` `panel.NewSnapshotPump(panel.PumpSources{...})` → `internal/panel/pump.go:207 Snapshot()` → `composer.go:106 NewSnapshot(...)`。r2/r3 记的 `run.go:442` 也是旧行号，现量 `:699-726`。`PumpSources` 十枚读口（`pump.go:121-188`）在装配根**全部接线**（`run.go:700-725`：Verdicts/Mode/Workspace/Git/Model/Credential/Results/Instructions/Tasks/Out）。

### 1A · `Snapshot` 直接字段（6 枚，`composer.go:57-92`）

| 字段 | json 键 | 真源（逐跳到行） | 态 |
|---|---|---|---|
| `Pending []ApprovalCardView` | `pending` | `pump.go:211-218` ← `Verdicts` 读口 ← `cmd/wisp/panel_pump.go:58 liveVerdicts` ← `rt.gate.Queue().LiveApprovals()` | **有源、已接** |
| `Results []ResultChunk` | `results` | `pump.go:220-223` ← `Results` 读口 ← `rt.stream.Chunks`（`pump.go:522` StreamLog） | **有源、已接** |
| `Composer ComposerState` | `composer` | `pump.go:242` `NewComposerState(...)`＋段内逐维覆盖 | **有源、已接**（见 1B，其中 2 枚恒空） |
| `GeneratedAt string` | `generatedAt` | `composer.go:131` `now.UTC()`（`run.go` 未注 `Now` 读口 ⇒ `pump.go:272` 取 `time.Now`）；仅显示、不推状态（`panel.ts:132-134` 语义） | **有源、已接** |
| `Instructions *InstructionsSection` | `instructions`（omitempty） | `pump.go:277-282` ← `Instructions` 读口 ← `run.go:718 rt.instructionBundle`（`panel_pump.go:108` 读 `instrLoader.Last()`） | **有源、已接**；nil 读口＝键缺席（诚实态） |
| `Tasks *TaskRosterSection` | `tasks`（omitempty） | `pump.go:283-291` ← `Tasks` 读口 ← `run.go:724 rt.taskRosterState`（`panel_pump.go:146` 走 roster） | **有源、已接** |

⇒ **顶层六枚全部有生产者**：票面 r1 立的第四态「无生产者」（carrier 在、无人读）对顶层已解除（200-r2 / 197-r3 / 248 各接了一根读口）。

### 1B · `ComposerState` 字段（11 枚，`composer.go:235-270`）

| 字段 | json 键 | 真源 | 态 |
|---|---|---|---|
| `Mode ModeView` | `mode` | `pump.go:229-247` ← `Mode` 读口←`run.go:701 rt.modes.PermissionMode`；不可读→`ModeUnknownView`（`composer.go:168`） | **已接**＋未知守卫 |
| `Workspace WorkspaceView` | `workspace` | `pump.go:234-237` ← `run.go:702 rt.workspaceView`（`panel_pump.go:83`→C26 canonicalizer） | **已接** |
| `Attachments []AttachmentRef` | `attachments` | `pump.go:242` 硬写 `nil`→`NewComposerState` 归 `[]AttachmentRef{}`；**无读口接线**，且 `AttachmentBroker`（`attachments.go:179`）无"列当前已挂附件"端口 | **①没接（恒空）**；真源存疑（附件随发出消息即消费，无持久"暂存列表"可量到） |
| `AttachmentMIMEs []string` | `acceptedAttachmentMimes` | `composer.go:281/119-120` `AcceptedMIMETypes()`（常量表） | **已接**（常量即真值，非零冒充） |
| `MaxAttachmentB int64` | `maxAttachmentBytes` | `pump.go:238-241` `AttachmentMax` 或默认 `MaxAttachmentBytes` | **已接** |
| `AttachmentReason string` | `attachmentError` | **泵从不填**→恒 `""`；无读口 | **①没接（恒空）** |
| `Git GitView` | `git` | `pump.go:248-252` ← `run.go:703 rt.gitView`（`panel_pump.go:229`→`ReadGitForWorkspace`）；不可读→`GitViewNotProbed`（`git.go:150`，Kind=`unreadable`） | **已接**＋未探守卫 |
| `CurrentModel string` | `currentModel` | `pump.go:257` ← `run.go:704 rt.currentModel`（`panel_pump.go:246`→`rt.endpoint.Model`，源 `internal/llm/resolver.go:46`、赋值 `run.go:308`） | **已接**（145-r3 落的那 1 维） |
| `ModelKnown bool` | `modelKnown` | `pump.go:258` `CurrentModel != ""` | **已接**（未知守卫伴位） |
| `Credential CredentialState` | `credentialState` | `pump.go:260-269` ← `run.go:710 rt.settings.credentialStatus`；五态词表 `config_handlers.go:177-191` | **已接**（248 落） |
| `CredentialKnown bool` | `credentialKnown` | `pump.go:269` 读口存在即置 true | **已接**（未知守卫伴位） |

### 1C · 段内其它字段（卡片与结果）

- `ApprovalCardView`（10 枚，`approval.go:39-59`）：correlationId／tool／args／level／rulesHit／reason／reasonKnown／sessionOverrideBlocked／decidedBy 均经 `CardViewFromDecision`（`approval.go:72`）从 live verdict 填＝**已接**；**唯 `callChain`＝③真源不存在**——`liveVerdicts`（`panel_pump.go:65-73`）不置 `NativeVerdict.CallChain`，且 `pump.go:80-86` 与 `panel_pump.go:50-57` 具名："run path records none，`tools.Decision` has no chain field"，只有 CLI 探针填自己的 argv 门。这是那张安全卡"欠用户链路却无源"的诚实空。
- `ResultChunk`（3 枚，`composer.go:95-99`）：correlationId（＝task id 或 `subagent:` 键）／text／done 全部由 StreamLog 填＝**已接**。**无 token 计数、无 reasoning**（那是 §3 的行 2/3）。
- `InstructionsSection`（`instructions_200.go:71-80`）：status/reason/files，五态互斥（loaded/none_found/off/refused/not_run，`:25-41`），非 loaded 必带 reason（`:158-193`）＝**已接**，且是"未知 vs 零"的样板。
- `TaskRosterSection`/`TaskRowView`（`subagent_roster_197.go:101-164`）：rows（含 `status`＋`statusKnown`＋`statusReason`＋`streamKey`＋`blockedOnApproval`＋每行截断三元）／inFlightSlots／poolCap／log 级截断三元，全部 `panel_pump.go:146 taskRosterState` 读 roster＋StreamLog＝**已接**；status 由 `tools.TaskOutput.StateAnswer`（宿主自判）带出，panel 侧无写腿（`TestTaskState188AC3PanelHasNoWriteLeg` 钉）。

### 1D · 快照里根本没有的维度（＝§3 会逐行归口；源在三家，跨不过 seam）

以下维度**无顶层/段内字段**，其源状态（本程现量）：
- 实时/累计 token 计数、成本金额：源＝`internal/agent/cost.go:23 Cost{Usage llm.Usage; Micros; Currency}`＋`guard.go:80-84` 的 `tokensIn/tokensOut/CostMicros`＋`llm.Usage`（`internal/llm/events.go:90`）——**真源存在**，但 `Cost`/`Guard` 是 `loop.go:367` 每次 run 的**局部对象**（`Loop` 结构 `loop.go:176-190` 只存 `guard GuardConfig` 配置、不存活体），到 `finish/brakeStuck` 才折进结果并被打印；**cmd/wisp 的泵读口无一经向 rt 拿到活体** ⇒ 属 ②"源存在、卡在 seam 之外"，非 ①（字段都没写）、非 ③（源真不存在）。
- reasoning 正文：源＝`llm.EvReasoningDelta`（`events.go:33`）三适配器真发、`events.go:189/258` 收集 `Reasoning` 字段；**到 run.go 记录点被打印未入快照**＝②。
- Stuck 重复计数＋工具名：源＝`guard.go:165 ObserveTurn` 产 `Reminder{Repeat,Tool,Stuck}`（`guard.go:47-58`），阶梯 `[3,5,8]`（`guard.go:97`）；同样是 loop 局部、不入快照＝②。**注**：Stuck/取消的**状态词**可能已折进 task 的 D43 filed status、经 `Tasks[].status/statusReason` 带出（`loop.go:378 finish(...StatusCancelled,"任务已取消")`、`brakeStuck`），此路本程**未量 StateAnswer 是否含该二词**（见 §6）。
- 错误人话文案 `humanText`：`observe` 有 Class/Detail（`guard.go:271 ErrorClassOfTurnError`→`observe.ClassOf`），但 **class→中文文案映射不存在**＝③。
- 注入命中片段 `fragment`：源形状在 `risk.Provenance` 的 `RiskHit.Fragment`（`provenance.go:278`，C25 取证用），但**不在 `risk.Decision`（`assessor.go:90-97` 四字段：Level/RulesHit/Reason/SessionOverrideBlocked）**⇒ 卡片拿不到＝**③对卡片路径**（源在别处存在，但未汇入 verdict→card；是否可在渲染时回读＝§6 量不到）。
- L1 环形倒计时 `remainingMs`：`approval.EventTick`（`ui.go:64`）**声明零枚发射者**；`Remaining` 字段（`gate.go:302/431`）带的是**静态窗口长 `g.window`**（`Gate.Window()` `gate.go:170` 可读），倒计时本体是 `gate.go` 的 `clock.After` 不暴露"还剩多少" ⇒ **实时 remaining＝③**；静态窗口长度＝有源可接但会"停在满格"（假相邻，r2/r3 拒收之因）。
- 思考/推理"已耗时 X.Xs"（`thinkingMs`/`reasoningMs`/`durationMs`）：**无计时生产者**（本程未找到任何 turn-started 时戳读口）＝③。
- 可选模型清单 `models`／各家档位词表 `efforts`：目录源在 `config`（`schema.go` ModelSpec∩Enabled，`run.go:209 rt.cfg` 有数据），但"∩enabled 列给人看"零枚生产读者（r2 P5 现量）；efforts 各家接受表全是**包内私有 `var`、无导出访问器**（r2 P8/§3.2：openairesponses 三档、anthropic 三档、openaichat 零映射）＝②（有源无读者）／efforts 近 ③（不可导出）。r3 明确排除＝票 187 射程。
- "当前屏" `view`／九枚屏 id：Go 侧无屏名常量表（r2/r3；`Q-51` 未答）＝③。

## §2 零值冒充真值的名册＋"未知 vs 零"的形状候选

**结论先给**：本仓对"宁缺毋造"已相当自律——**§1 里后来新增的每一维都带着一枚"未知守卫"**，没有发现"把没测过画成空闲/零"的现役字段。真正恒空的字段是 `attachments`/`attachmentError` 两枚（无真源/无读口），它们的空是"确实没有附件"意义上的诚实空，但**缺 known 位**，于是"看不见附件"与"没有附件"在页面上塌成同一读——这正是票面要防的那一形。下面分"已有的三种守卫形状"＋"残余风险名册"＋"最少形状候选"三段。

### 2A · 本仓已落的三种"区分未知与零"形状（现量 file:line）

| 形状 | 用在哪枚字段 | 赋值行 | 说明 |
|---|---|---|---|
| **伴生 `xxxKnown bool`** | `ModelKnown`/`CredentialKnown`（`composer.go:256,269`）、`TaskRowView.StatusKnown`（`subagent_roster_197.go:110`）、`ApprovalCardView.ReasonKnown`（`approval.go:52`） | `pump.go:258/269`（读口存在才置 true）、`subagent_roster_197.go:215-223`（StatusKnown=false 时把 Status 抹成 ""＋填 fail-closed reason） | 标量维的标准做法：值＋一枚"我到底读没读到"。计数类维度（0 是合法值）**只能**用这一形，没有天然 sentinel。 |
| **保留 sentinel 枚举成员** | `ModeView.Current="unknown"`、`CredentialState="unknown"/"config_unreadable"`、`GitView.Kind="unreadable"`、`InstructionsSection.Status="not_run"/"none_found"/"off"/"refused"` | `composer.go:117`、`config_handlers.go:180,183`、`git.go:151 GitViewNotProbed`、`instructions_200.go:159-164` | 字符串/枚举维的标准做法：给"没读数"留一个**具名**档位，而不是留空串。 |
| **非真态强制 `Reason` 非空** | `GitView.Reason`（`git.go:120`"mandatory, never empty"）、`WorkspaceView.Reason`（`composer.go:222-224` Set=false 时必带）、`InstructionsSection.Reason`、`TaskRowView.StatusReason` | 上述各构造点 | 前两者再配一句"为什么没有"，杜绝渲染侧自行补"没有风险/没有仓库"。 |
| **段级 `*`＋omitempty＝键存在性表意** | `Snapshot.Instructions`/`Snapshot.Tasks`（`composer.go:74,91`） | `pump.go:277-282/283-291`（读口在才挂段；读口在但空内容仍发键，读口无才让键缺席） | 整段级："这个 host 根本看不见子代理/没接加载器"靠**键缺席**表达，而非靠一段全零——`composer.go:66-73` 明写这不是绕 byte-nail 的近路。 |

⇒ **顶层六段＋composer 内 9/11 维都已被上述某种形状守住了**；构造器 `NewSnapshot` 自己还兜了两道（`composer.go:115-126`：mode 空→unknown、git 空→NotProbed）。

### 2B · 残余"零值冒充真值"风险名册（具名到行）

1. **`ComposerState.AttachmentReason`（json `attachmentError`，`composer.go:241`）——无 known 位的恒空**。泵从不填（`pump.go:242` 只传 nil，段内无一处写它），页面读到的一直是 `""`＝"没有附件错误"。风险等级低（因为 `Attachments` 也恒空，两者自洽），但它**没走 2A 任一种形状**：将来若接上"有读口但本轮读失败"，`""` 会把"读失败"画成"没错误"。**这是本程唯一发现的、现役字段缺未知守卫处。**
2. **`StreamLog.TruncationFor` 未知键→零值**（`pump.go:655-668`：`s.kept[key]` 不存在即 `return StreamTruncation{}`）。于是 `TaskRowView` 的 `streamTruncated=false / streamElidedRunes=0 / streamDropped=false`（`subagent_roster_197.go:130-132`、填于 `panel_pump.go:172,182-184`）对"这条流从没写过"与"这条流完整没截"给同一读数。`pump.go:651-655` 的注释**自知并接受了这一塌合**（"the row that exists is the row in full"）；不算新 bug，但属"零冒充真"的既登记缺口，页面若把 `elided=0` 画成"已确认全文"就越了界。
3. **假想的 `remainingMs`**（当前快照**没有**这枚字段）：若某写腿把 `gate.go:302/431` 的静态 `Remaining: g.window` 直接搬进快照当"还剩多少"，就会**把"刚开满格窗口"画成"还在倒数"**——`pump.go:50-56` 具名此坑（`EventTick` 零发射、无"还剩多少"读口）。**现状＝字段不存在，所以没违反；风险在下一步 naiive 接法**。r2/r3 正是因此拒收 `remainingMs`。
4. **`Cost.Micros` / `Guard.tokensIn+Out` 若接进快照而不带 known 位**：金额/计数 0 是合法值（任务确实没花钱），"没读到 Cost（loop 局部，泵拿不到）"与"确实 0"必须分开——接的时候**必须**用 2A 第一形（伴生 known bool），否则就是把"无源"画成"零成本"。

### 2C · "未知 vs 零"最少需要的形状候选（只列形状与代价，⛔ 不选形、不动手）

给编排者排形用，按维类型分列，代价＝对**契约/前端对齐/泵**三处的连带成本：

- **标量字符串/枚举维**（如 model、mode、credential 一族）→ **形状＝保留 sentinel 成员**（2A 第二形）。代价：新增一个枚举名（前端 TS 要能识别；`composer_test.go:48` 双向尺只认 json 键不认值， sentinel 不触发它），零新增字段。
- **标量计数/布尔维**（0/"" 是合法真值，如 cost、elided、remaining）→ **形状＝伴生 `xxxKnown bool`**（2A 第一形）。代价：**+1 字段/维＝+1 个 json 键＝扩契约**，会撞 `approval_test.go:105`/`composer_test.go:48` 双向尺与 `pump_test.go` 键集钉（须同 commit 带 `frontend/src/lib/panel.ts`＝`Q-51` 未决）；泵侧多一条"读口在才置 true"的赋值线。
- **整段/整块维**（如 tasks、instructions、未来的 tools[]/failures[]）→ **形状＝`*Section`＋omitempty，键存在性表意**（2A 第四形）。代价：段对象一个 json 键（比逐维 known 少扩），但要保住 4 键地板的 byte-nail（`composer.go:66-73` 已论证指针 omitempty 不违反），且前端要为"键缺席"与"键在但空数组"写两套文案。
- **最少之我见（不选，只摆）**：一维若"没读到"与"值本身是零"都可能出现 → **非要有某种 unknown 表示，三形择一**；纯计数维**排除 sentinel 形**（没有天然"未知"数值可借），只 `known bool` 或 `*int`（nil＝未知）二者之一——但 `*int` 也是扩一个可空键、代价同上。

> ⚠ 本节仅**列形状＋代价**，选形归编排者；本程未动任何产码、未提议放宽任何既有断言（键集钉 `pump_test.go:111-124/:270-276`、双向尺 `composer_test.go:48`/`approval_test.go:105` 一律照其在位）。

## §3 状态↔缺字段对照表（十四枚逐行）

〔待填〕

## §4 复认/推翻前人件（具名）

〔待填〕

## §5 我可能写错的条目（自我对抗）

〔待填〕

## §6 量不到的地方（具名）

〔待填〕

## §7 交件判语

〔待填〕
