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

**十四行来源＝`docs/PLAN.md:3481-3494`（本程现读；票面引的 `:3473-3488` 已漂，表头在 `:3479`，行体从 `:3481` 起——见 §5）。**
**关键再框定**：本票标题说"快照只 4 字段→11 态无输入"。现量快照确实长到 6＋11 字段，但**新长的维度（git/model/credential/instructions/tasks）几乎没有一枚喂这十四行"工作态"**——它们是"设置/名册"侧的维，不是"这一轮 agent 正在干什么"的侧。⇒ 十四行里的**工作态**基本**仍未接**；下表逐行标"今天能否画"。

| # | 状态（PLAN 行） | 今天能画？ | 缺的字段（Go 侧候选名） | 该字段属哪种"没做" | 盘上证据 |
|---|---|---|---|---|---|
| 1 | 思考中（等 LLM 首 token） | **否** | `run.phase`（"正在思考"这一相）＋ `thinkingElapsedMs`（已耗时） | phase＝**③**（无 phase/状态映射表，r2/r3）；elapsed＝**③**（无 turn-started 时戳读口） | 快照无 phase 字段；`grep` 计时读口无果（§6 待量） |
| 2 | SSE 流式输出 | **部分**（文字有，实时 token 无） | `results[].tokens`（实时 token 计数） | **②**（源＝`agent.Cost.Usage`/`llm.Usage` `cost.go:24`、`events.go:90`，卡在 loop 局部跨不到泵） | 文字：`composer.go:97 Results[].text`（已接）；计数：无字段 |
| 3 | 推理过程（reasoning delta） | **否** | `run.reasoning`（正文）＋ `reasoningElapsedMs` | 正文＝**②**（源＝`llm.EvReasoningDelta` `events.go:33`，收集器 `events.go:189,258 Reasoning`；记录点被打印未入快照）；elapsed＝**③** | r2 §2 `run.reasoning` 行：源真、卡在 `run.go` 记录点 |
| 4 | 工具调用（chip） | **仅被 L2 门控的那几枚**（=卡片）；运行中/非门控工具调用**无法画** | `tools[]`（条目＋四值状态 running/success/fail/rejected＋风险徽标） | **②**（条目源＝record point `run.go:816` 一带 / `agent` 事件；四值状态无快照承载；非门控调用完全无读口） | 门控调用＝`pending[].ApprovalCardView`（level/tool/args 有，`approval.go:39-59`）；其余无 |
| 5 | 工具调用（展开） | **部分**（门控调用有全参数，结果摘要/耗时无） | `tools[].result`＋`tools[].durationMs`＋`correlationId`（已有） | result＝**②**；durationMs＝**③**（无计时生产者） | 全参数＝`approval.go:42 args`；correlationId=`approval.go:40`；result/duration 无字段 |
| 6 | 审批等待（L2 队列深度角标） | **能画**（`depth = len(pending)` 客户端可算） | 无（**不需新字段**：r2 把 `approval.depth` 判为"复述字段当装饰"拒收） | —（可派生） | `composer.go:58 Pending`；`panel_pump.go:333-339` 摘要已用 `len(snap.Pending)` |
| 7 | L2 确认卡（面板内） | **能画**（唯一"已有"的一行，除 callChain 外齐） | `callChain`（链路）＝缺 | **③**（`pump.go:80-86`/`panel_pump.go:50-57`：run 路径不记链路，`tools.Decision` 无 chain 字段） | `approval.go:39-59` 九枚已接；`approval.go:55 callChain` 恒空 |
| 8 | L1 阻止窗口 | **部分/否**（工具名与参数摘要视是否进 LiveApprovals 而定；环形倒计时**无法画**） | `approval.remainingMs`（实时还剩多少） | **③**（`EventTick` 零发射 `ui.go:64`；`gate.go:302/431 Remaining` 带的是静态 `g.window` 非实时剩余） | `pump.go:50-56` 具名此坑；静态窗长可经 `Gate.Window()` `gate.go:170` 取但会"停在满格" |
| 9 | L2 原生降级卡（C27，面板不可用时） | **面板快照射程外**（原生球侧弹出，非 Go→页面快照） | —（非本载体） | —（票面①：原生侧、面板契约够不着） | 归 C27 原生路径，非 `Snapshot` |
| 10 | 错误 | **否** | `failures[]`（class＋detail）＋`failures[].humanText`（人话） | 载体/class＝**②**（`observe.Error` Class/Detail `guard.go:271 ErrorClassOfTurnError`；事件在 record point，`agent.Event.Err` 真发不留存）；humanText＝**③**（class→中文映射不存在，r2 §2） | 无任何 error 字段 |
| 11 | Stuck（梯度刹车 C22） | **待量**（若 `Tasks[].status/statusReason` 已带 Stuck 词则可画；否则否） | `stuck`（＋重复的工具名与次数） | 源＝**②**（`agent.Guard.Reminder{Repeat,Tool,Stuck}` `guard.go:47-58,165`，阶梯 `[3,5,8]` `guard.go:97`；loop 局部、未入快照）；是否已折进 task status＝§6 待量 | `loop.go:382-388 brakeStuck`；快照无独立 stuck 字段 |
| 12 | 成本（C23） | **否** | `cost`（IN/OUT/CACHED token＋金额） | **②**（源＝`agent.Cost{Usage,Micros,Currency}` `cost.go:23`、`guard.go:80-84`）；⚠ 单位口径 micro-USD vs CNY **未定案**（`cost.go:16-20`"never invents exchange rate"）⇒ 属**待裁**非纯字段 | 快照无 cost 字段；r2 §2 拒收（单位＋跨 seam 双因） |
| 13 | 已取消/被打断 | **待量**（取消原因串源在 `loop.go:378 finish(...StatusCancelled,"任务已取消")`→可能经 `Tasks[].status/statusReason` 带出；否则否） | `cancelReason`（若未经 tasks 带出） | 原因串＝源存在（`queue.go:501 abandon(why)`、`loop.go:378`），是否达快照＝§6 待量；若不达＝**②** | 无顶层 cancel 字段 |
| 14 | 注入检出（C25/R4） | **部分**（"含来自 `<url>` 的内容"这句 reason **能画**；"可展开看命中的片段"**不能**） | `card.fragment`（命中片段） | **③ 对卡片路径**（`RiskHit.Fragment` `provenance.go:278` 有源，但**不在 `risk.Decision`** `assessor.go:90-97`→不到卡；渲染时可否回读＝§6） | 票面④：来源链 rules_gateway→`approval.go:83 Reason`→卡 **是通的**；缺陷只是硬编码文案（前端自改），非等 Go |

**小结（本程现量口径）**：十四行里 **今天能画的＝行 6、行 7（除 callChain）** 两行完全、**部分能画＝行 2/4/5/14** 四行、**待量（可能经 Tasks 带出）＝行 11/13** 两行、**完全不能画＝行 1/3/8/10/12** 五行、**射程外＝行 9（原生）**。
⇒ **对"工作态"而言票面①的"缺 8～10"判断方向仍成立**：快照新增的 6 字段（git/model/credential/instructions/tasks）解决的是"设置/名册"侧可读性，**没有一枚落进这十四行的工作态缺口**。真正的跨-seam 阻塞仍是那句"数据在 `internal/agent` 的 loop 局部/记录点算得出、泵读不到"（行 2/3/10/11/12）＋"纯无源"（计时 elapsed、人话文案、命中片段、实时剩余、phase 映射）。

## §4 复认/推翻前人件（具名）

前人件：①`docs/evidence/s1/145-snapshot-field-census-r1.md`（AC#1 普查，锚 `88eab34`）＋票面 ①–⑦ 更正；②`145-snapshot-growth-r2.md`（14 候选逐枚）；③`145-snapshot-growth-r3.md`（落 1 维）。本程对 HEAD `5f9ff9d4` 现读复认。**注意：多处前人读数并非"错"，而是被 r1 之后落地的票（181/197/200/248）与票 35 装配根"跑在了前面"——本程区分"错"与"过期"。**

### 4A · 复认（与现量一致）

1. **票面① / r1 P7："十一行"不成立，面板侧完全无输入＝8（＋行 8 半＝9，＋行 9 原生＝10）。** ⇒ **复认**，且本程 §3 给出更强理由：快照虽长到 6 字段，**新增维度（git/model/credential/instructions/tasks）无一喂这十四行工作态**，"缺"的方向仍成立。
2. **r3：`composer.currentModel`＋`modelKnown` 已落、真源＝`llm.Endpoint.Model`。** ⇒ **复认**：`composer.go:255-256`、`pump.go:253-259`、`run.go:704`、`panel_pump.go:242-247`、源 `internal/llm/resolver.go:46`＋赋值 `run.go:308`。
3. **r1/r3：`ApprovalCardView.callChain` 无生产者（run 路径不记链路）。** ⇒ **复认**：`approval.go:55` 字段在、`panel_pump.go:65-73 liveVerdicts` 不置 CallChain、`pump.go:80-86` 与 `panel_pump.go:50-57` 具名。
4. **r1 §2.4（行 4）"工具运行中不是无源，是接缝有、转发无"。** ⇒ **复认**：`agent.EvToolStart` 声明于 `internal/agent/sink.go:25`，全仓**零枚发射点**（本程 grep 只命中声明＋`run.go` 消费分支）；而底层 `llm.EvToolCallStart` 真发（`internal/llm/anthropic/stream.go:262`、`events.go:34`）。⇒ 工具条目属 ②（源在 llm 事件，缺 agent 级转发＋快照承载）。
5. **r2/r3：`approval.remainingMs` 实时剩余无生产者。** ⇒ **复认并现量钉死**：`approval.EventTick` 声明于 `internal/agent/approval/ui.go:64`、**零发射者**；`gate.go:302/431` 的 `Remaining: g.window` 是**静态窗口全长**（`Gate.Window()` `gate.go:170`），非"还剩多少"。
6. **票面④（行 14）：注入来源链今天通（非"送不到面板"）。** ⇒ **复认**：卡片 `Reason` 一路带过来（`approval.go:82`，`rules_gateway→gate→card`），"含来自 `<url>` 的内容"今天能画；缺的只是**可展开的命中片段**（见 4C）。

### 4B · 推翻 / 需修正（现量与旧读冲突）

7. **r1 §2.4（行 8）判 L1 `RemainingMs/Channels/WindowMs` "全有源"——本程判其把"静态窗长"与"实时剩余"混为一谈。** ⇒ **推翻这一枚"全有源"**（与 r2/r3 一致、与现量一致）：WindowMs（窗长）确可读，但行 8 视觉要的是**环形递减**＝实时 remaining，而 EventTick 零发射 ⇒ 那枚具体读数无生产者。**注**：`observe.Timeout.Remaining()`（`internal/observe/clock.go:43`）是一枚现成的单调剩余原语，`agent/tools.go:180` 已用它做工具超时——**L1 若开一枚 Timeout 记窗起即可派"还剩多少"**，所以此行不是"永远无源"，是"尚未为 L1 建实例＋接线"（比 r2/r3 的纯"无源"更精确）。
8. **r2 进度栏（票内 151–153）"真正没源的仍是那 7 枚"一句，把 cost/usage/reasoning/repeat 也扫进"无源"——本程推翻其措辞、收窄其范围。** ⇒ 现量：**cost/usage 有活源**（`internal/agent/cost.go:23` `Cost{Usage,Micros,Currency}`＋`guard.go:80-84` `tokensIn/tokensOut/CostMicros`＋`llm.Usage` `events.go:90`）；**reasoning 正文有源**（`llm.EvReasoningDelta` `events.go:33`＋收集器 `events.go:189,258`）；**repeat/stuck 有源**（`guard.go:165 ObserveTurn`→`Reminder{Repeat,Tool,Stuck}` `guard.go:47-58`）。它们的阻塞是**跨 seam**（这些都活在 `internal/agent/loop.go:367` 每次 run 的**局部** `guard`/`cost` 上，`Loop` 结构 `loop.go:176-190` 只存配置不存活体，cmd/wisp 泵读口无一触及）＝ **② 非 ③**。⚠ 但 r2 的**详表**（§2/§3.2/P5/P7/P8）其实自己就写了 models/capabilities 是"有源无读者/半有源"——所以这是**r2 进度栏一行话与 r2 详表内部的张力**，本程站详表那侧。
9. **票面⑤ / r1 §3(1)："今天没有任何一处生产代码构造快照"（`NewSnapshot` 非测试调用者＝0）。** ⇒ **过期推翻**（非 r1 之错）：票 35 装配根已落地——`pump.go:276` 调 `NewSnapshot`，`run.go:699` 调 `NewSnapshotPump` 并接满十枚读口，`panel_pump.go:401 publishPanelSnapshot` 在状态移动时驱动。r1 立的第四态"无生产者"对**已落地的那些维**（pending/results/composer 各维/instructions/tasks）**已解除**。
10. **票面标题 / r1 P1："`Snapshot` 恰四字段"（`composer.go:44-49`）。** ⇒ **过期推翻（针对当前 HEAD）**：现量 6 枚直接字段（`composer.go:57-92`）。根因＝Instructions(200)、Tasks(197) 在 r1 锚 `88eab34` 之后加进结构体；r1 当时"四枚"成立。⚠ **连 `composer.go:40-56` 自己的头注释都还写着"FOUR, not one"并引装配根"`run.go:421-427`"**（现 `:699`）——这是一处**文档漂移**（本程只读、不改，据"薄索引≠权威"规约上报）。
11. **前人所有行号相对 HEAD 均漂。** ⇒ 本程现读为准：`Snapshot` type `:57`（字段 `:58-61`＋`:74`＋`:91`）、`NewSnapshot` `:106`（return `:127`）、`ComposerState` `:235-270`、装配根 `run.go:699-726`、十四行表头 `PLAN.md:3479`／行体 `:3481-3494`。r3 引的 `run.go:442` 现量＝旧行（现 `:699`），`run.go:800-816`"记录点"现量已非记录点（`:801-813` 是 attachReplyListener/startConfigReload）。

### 4C · 一处本程**未复认**的否证（不据此下结论）

- **r1 §2.6（行 5，`DurationMs`）"tool_call.started_at 在生产路径上没有任何写者"**——本程**未复认**这枚否证（不跑码、未读到 record-point 的 `TaskLog.ToolCall.StartedAt` 赋值处）。现量只到：DB 列与 DAO 齐（`internal/memory/dao_toolcall.go:25,29,174`、`models.go:60 StartedAt int64`、`dao_tasklog.go:27,33` task 级 started_at **确有写者**）。⇒ "工具级 duration 无写者"是 r1 的否证、本程既没坐实也没推翻，落 §6 待量；**下游勿据此否证安静绕路**。

## §5 我可能写错的条目（自我对抗）

〔待填〕

## §6 量不到的地方（具名）

〔待填〕

## §7 交件判语

〔待填〕
