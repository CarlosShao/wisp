# 182 — 右上角「任务监控」栏 10 堆的 Go 侧源、载体与归属票（r2 重跑；锚 `5183f9715` 起手 → `a42d6fccc` 收尾）

> 票 182，AC#2／AC#3／AC#4／AC#5 的可判证据（本腿**不翻任何框**，勾归编排者）。
> **只读腿**：`internal/**`、`cmd/**`、`docs/PLAN.md`、`docs/specs/**`、`design/**`、`frontend/**` **零字节写入**；
> 零产码；写面只有本文件（新建）＋票 182 Progress log 一行追加。
> 本件是 `docs/evidence/s1/182-census.md`（下称 r1，锚 `232b0846c`）的**重跑与订正**，不是它的复述：
> r1 的读数在今天的树上有 **6 处过期**（§7 逐条具名，每条附尺），本件以今天的工作树为准。
> **口径**：堆数沿用编排者 09-28 17:5x 的裁定＝原型 `design/doubao/demo/**` 注册表口径，**不是产品真身**。
> owner 那张截图是**外部产品（qoder）的界面，不是本仓规格**（AC#1 要求明写；本件没有一条读数来自截图）。
> **正控纪律**：凡 0 命中，同一次调用里都打了「已知存在的同形尺」，读数记在括号里；没打正控的不记数。

---

## 1. AC#1 复核：堆数 = 10（7 枚有位 + 3 枚有名无位）

| 尺 | 原文 | 今天读数 |
|---|---|---|
| 主尺 | `grep -rn --include=*.js "RbPanels.register(" design/doubao/demo/` | **7 行** |
| 对照尺 | `grep -rn --include=*.js -E "^\s*label:" design/doubao/demo/` | **8 行**，剔 `app.js:300`（命令面板条目，不是竖条堆）＝**7**，与主尺相等 |

**7 枚有位**（`label` 逐字）：

| 位 | 堆 | 注册位 | demo 里的子行 |
|---|---|---|---|
| K1 | 上下文 | `rightbar.js:8` | CWD／附件两行／「上下文预算 42%」＋IN/OUT/CACHED |
| K2 | 活动 | `rightbar.js:83` | 活动时间线（4 条：fs.listdir 完成／task.extract 运行中／fs.delete 被拒 L2／审批超时自动拒绝） |
| K3 | 审批 | `rightbar.js:116` | 待审批队列 9/10＋满额告警＋条目（只跳屏，无允许/拒绝按钮） |
| K4 | 文件 | `rb-files.js:6` | CWD 面包屑＋工作区文件树＋路径锁说明 |
| K5 | 插件 | `rb-plugins.js:6` | 技能／伙伴两组＋`on:` 开关＋`level:` |
| K6 | 审阅 | `rb-review.js:6` | 逐文件 diff＋「Apply 3 changes」 |
| K7 | 终端 | `rb-terminal.js:6` | 标签页（main/build）＋「新建终端」＋命令输入行 |

**3 枚有名无位**（今天不在注册表；仓内文本有据）：

| 位 | 堆 | 仓内出处（本程现读） | demo 里有位吗 |
|---|---|---|---|
| K8 | 子代理状态 | 票 188 AC#1（owner 已裁「乙：进」）＋票 197；`issues/77:364`「右栏三标签定案：审查/终端/浏览器」同段 | 无（`grep -rn "子代理" design/doubao/demo --include=*.js`＝**0**） |
| K9 | 后台任务名册 | 票 188 AC#2＋票 164（已结案）＋票 196（已结案） | 无（demo 的 K2 是工具活动流，不是任务名册） |
| K10 | 浏览器 | `issues/77:364` 逐字「右栏三标签定案：审查/终端/浏览器」＋`:384` 逐字「右栏 审查/终端/浏览器三面板、任务监控弹窗」（两行本程现读对位） | 无（三标签里那一枚今天没有文件） |

⇒ **可点名总堆数 = 10**，与 r1 同数、同成员。本程**不翻 AC#1**（10-08 已由编排者翻）。

---

## 2. AC#2 逐堆三档定性（每档附尺与命中数）

三档定义（票面原文）：**有源** = Go 里已能取到真值（缺不缺载体另说）／**无源但规格要求** = Go 取不到但 `PLAN.md`／specs 有人要求过／**规格真空** = 今天没人要求过，做 = 新需求。

### 2.1 总表

| 堆 | 三档 | 现读尺（原文）→ 命中数 | 今天载体在哪 | 要动哪一节／是否契约 |
|---|---|---|---|---|
| **K1-A 路径（CWD）** | **有源＋载体已在** | `grep -n "Workspace .*json:\"workspace\"" internal/panel/composer.go`＝**1**（`:273`）；`awk '/^type ComposerState struct/,/^}/' internal/panel/composer.go \| grep -c "json:"`＝**12** | 快照 `composer.workspace`；装配 `Workspace: rt.workspaceView`（`cmd/wisp/run.go:930`） | 不动（票 92 已 `-done`） |
| **K1-B 附件** | **有源＋载体已在** | `composer.go:274-277` 四枚（`attachments`／`acceptedAttachmentMimes`／`maxAttachmentBytes`／`attachmentError`） | 快照 `composer.*` | 不动（票 92） |
| **K1-C 占用条（used÷total）** | **两半源都在、口径未定、载体缺**（对 r1「无源」的收窄） | 分母候选＝`internal/agent/budgets.go:79 TokenBudget`＋`:50 Budgets.ContextWindow`（现读）；分子＝`Result.Usage`（`cost.go:52 AddUsage` 累加）；**面板侧**：`grep -rnE "usedTokens\|UsedTokens\|contextUsed\|ContextUsed\|Budget" internal/panel --include=*.go \| grep -v _test`＝**0**（正控：同包 `json:"`＝**191** 行） | 无；`PumpSources`（`internal/panel/pump.go:138-233`）11 个 reader 里没有占用 reader | 票 **167 AC#2**（未勾）：分子或分母任一未知 ⇒ 出口必须是「未知」这一枚具名状态；**口径未定 ⇒ 不许先落字段** |
| **K1-D token／成本** | **有源缺载体** | 源链：`internal/agent/cost.go:52 AddUsage` → `guard.go:148 AddUsage` → `loop.go:536`／`:540` → `:1115 publishUsage`（写 `Result.Usage/CostMicros/Currency`）；**面板消费**：`grep -rniE "\\bUsage\\b\|\\bCostView\\b\|tokensIn\|totalTokens" internal/panel --include=*.go \| grep -v _test`＝**0**（正控同上，同一把 191 行的 `json:` 尺在同一包里响） | **无**：`Snapshot` 6 键（`pending/results/composer/generatedAt/instructions/tasks`）里没有 cost；`PumpSources` 无 Cost reader | 票 **145 行 12「成本（C23）」**，今天位＝`PLAN.md:3510`（表 `:3499-3512`）；145 **AC#2 未勾**、AC#2b 未勾、AC#6 未勾；⚠ 单位口径（micro 单位 vs 币种）在 `cost.go:14-18` 自陈「open」，是 145 行 12 的 ⚠ 那一格 |
| **K2 活动** | **有源缺载体（细维）** | 源：`internal/agent/sink.go:53 type Event`（`Kind/TaskID/Text/ToolName/CallID/Outcome/Verb/Stop/TokensIn/TokensOut/RepeatLevel/Err`）；面板：`grep -rnE "json:\"feed\"\|ActivityView" internal/panel --include=*.go \| grep -v _test`＝**0**（正控：`json:"rows"`＝**1**，`subagent_roster_197.go:158`） | `Snapshot.results` 只有 `ResultChunk{correlationId,text,done}` 三枚（`composer.go:78-82`），**无工具名/状态/时刻** | 内容维**已有规格**：`PLAN.md:3502` 行 4「工具调用」chip 逐字定义了「运行中／成功／失败／被拒 ＋ L0/L1/L2」⇒ 并 **145 行 4／5 的载体扩张**；若要「独立时间线」＝十四行表**新增一行＝改 `PLAN.md` 文字（须人工批）** |
| **K3 审批** | **有源＋载体已在**（序号/停止两维缺） | `ApprovalCardView` **10** 个 json 字段（`internal/panel/approval.go`）；`grep -rn "Position" internal/panel --include=*.go \| grep -v _test`＝**0**（正控：`approval.go` 里 `CorrelationID`＝**3**）；真身：`internal/agent/approval/pending_read.go:58 Position`、`:127` 赋值 | 快照 `pending` | 载体格＝145 行 6／7；**序号＝167 AC#3、停止＝167 AC#4**（皆未勾）；不新增键 |
| **K4 文件树** | **有源＋载体已在**（Go 侧） | `TreeView` **7** 个 json 字段（`internal/panel/tree_types_190.go:50-58`）；`Files *TreeView`＝`composer.go:285`；装配 `Files: rt.fileTreeView`＝`run.go:932`；只读：`grep -nE "os\\.Remove\|os\\.Rename\|os\\.WriteFile\|os\\.Create\|RemoveAll" internal/panel/tree_reader_190.go`＝**0**（正控：同文件 `os.Open\|os.Lstat`＝**2**） | 快照 `composer.files`（`omitempty`） | 不动（票 190 已结案 6/6）；⚠ **跨树红在册**：`TestComposerContractTypesMatchFrontend` 今天红，报「Go ComposerState emits [files] that interface ComposerState does not declare」＝前端声明欠账（`Q-51` 那一族），**不是 Go 源缺** |
| **K5 插件** | **库里有读面／生产零装配／栏位规格真空** | `func (s *Store) ListPluginStates`＝`internal/memory/dao_misc.go:243`；**非测试调用者 0**（`grep -rn "ListPluginStates" internal cmd --include=*.go \| grep -v _test.go \| grep -v dao_misc.go`＝**0**；测试调用 2：`dao_test.go:451`／`manifestslot/install_50r1f_test.go:403`）；面板侧 `grep -rni "plugin" internal/panel --include=*.go \| grep -v _test`＝**3 行**，逐行过目＝1 行锁定族名单（`config_handlers.go:109`）＋2 行散文（`config_write_313.go:33`／`display_scrub_216.go:45`），**零读面** | 无 | **票 407 认领（2026-11-05 新建）**，但该票自陈 **Blocked by「前端 face 归外部 agent」** ⇒ 本编队写面只到「Go 侧 bridge 读口若缺另立子格」；**开关写面另见 §5 M3，不并进只读交付** |
| **K6 审阅（未提交枚数／逐文件 diff）** | **有源＋载体已在** | `GitView` **12** 个 json 字段（`internal/panel/git.go:144-181`，含 `Uncommitted`＝`:178`、`Diff`＝`:181`）；`GitDiffView` **9** 个 json（`git_readonly_exec.go:200-`：`kind/reason/files/filesTotal/filesIncluded/bytesIncluded/byteCap/...`）；`GitDiffFile` **10** 个 json（`:163-196`）；进程面 `runGitReadOnly(` 出现 **4** 次（定义 `:349` ＋ 调用 `:426`／`:493`／`:556`）；写动词字面 `grep -nE "\"(commit\|push\|checkout\|reset\|apply\|restore)\"" internal/panel/git_readonly_exec.go`＝**0**（正控：查询字面 `"status"\|"diff"\|"rev-parse"\|"ls-files"`＝**5**） | 快照 `composer.git.uncommitted` / `.diff` | 不动（票 189 已 `-done`）；分工见 §4 |
| **K7 终端** | **无源但规格要求** | 注册循环 `grep -c "tools.Builtin.*Entries(" cmd/wisp/run.go`＝**5**（`:624` FS／`:653` Task／`:680` System／`:1038` Subagent／`:1077` Plan，**无 Shell 一支**）；`grep -rn "shell\\.session\|shell\\.exec" internal cmd --include=*.go \| grep -v _test.go`＝**6 行全注释**（`config/unwired.go:63/64/70/76`、`risk/rules_shell.go:14`、`agent/spill.go:26`）；PTY 非测试 1 命中＝`cmd/wisp/testdata/esclistener/main.go:117` 的 `ptY` 假阳 | 无 | 契约行：D34 表 `PLAN.md:2575` 逐字「`shell.session`｜在一条活着的会话里跑下一条命令｜**L2**｜票 163」；地基＝**票 163 AC#1b**（未勾）；载体＝**票 191**（AC#0 已裁**甲**：复用工具腿＋门控＋审计） |
| **K8 子代理** | **有源＋载体已在** | `TaskRowView` **12** 个 json 字段（`internal/panel/subagent_roster_197.go:111-142`，含 `Status`／`StatusKnown`／`StatusReason`＝`:119-125`）；装配 `Tasks: rt.taskRosterState`＝`run.go:953` | 快照 `tasks` | 不动规格；票 **197**（2 勾／5 未勾）＋票 188 AC#1 已裁「乙：进」 |
| **K9 后台任务** | **有源＋载体已在** | `TaskRosterSection` **7** 个 json 字段（`:155-183`：`rows`／`inFlightSlots`／`poolCap`／三个流截断读数）；状态词表＝`tools.TaskOutput.StateAnswer` | 同上 `tasks` | 票 **188 AC#2**（未勾；⚠ 前提「给 `Record` 补状态维」与代码对不上，见 §9-2）＋票 164／196（皆 `-done`） |
| **K10 浏览器** | **无源＋规格真空（有票，被前置扣住）** | `grep -rli "browser" internal/panel --include=*.go \| grep -v _test`＝**0**（正控：`grep -rli "WebView" cmd/wisp --include=*.go \| grep -v _test`＝**14 文件**） | 无 | **票 192**（AC#1–#5 未勾；AC#0 已定案「零 key 抓结果页」）＋**票 168**（能力/插件门控；0/7 未勾）；分工见 §3 表脚注 |

### 2.2 小结（按「只差什么」排）

1. **只差载体（票 145 那一集）**：K1-D token／成本（行 12）、K2 活动（行 4/5 扩载体）。
2. **只差一个口径或一个字段维（票 167）**：K1-C 占用（AC#2，口径未定）、K3 序号（AC#3）、K3 停止（AC#4）。
3. **已闭、只须记账**：K1-A／K1-B（92）、K3 主体（`pending` 已是快照第一维）、K4（190）、K6（189）、K8／K9（197＋188 AC#1／164／196）。
4. **无源、票已在手但不许本票重开**：K7（163 地基未注册＋191 载体未落地）。
5. **认领已变但仍开**：K5 读＝407（Blocked by 前端 face）；K10＝192（＋168 能力侧）。

---

## 3. AC#3 逐堆归属票（本票主要产出）

> 指认前逐张票现读了 AC 标题与框态（读数在「框态」列）。**本程没有替任何票翻勾**，也没有在票面外新增任何票。

| 堆 | 归属（具名 + 格） | 并格还是立票 | 该票今天框态（现读） | 一格一句为什么是它 |
|---|---|---|---|---|
| K1-A 路径／K1-B 附件 | **票 92** | 已闭 | `92-...-done.md`；6 勾 | `workspace` 与附件四枚都是票 92 落的 |
| K1-C 占用条 | **票 167 AC#2** | 并格 | 1 勾／6 未勾 | 该格原文就是「分子或分母任一未知 ⇒ 出口里必须是『未知』这一枚具名状态」 |
| K1-D token／成本 | **票 145 行 12** | 并格（145 的「扩快照」） | 4 勾／3 未勾（AC#2／AC#2b／AC#6） | 行 12 逐字「成本（C23）」；源链 `cost.go`→`guard.go`→`loop.go` 现读可得 |
| K2 活动 | **票 145 行 4／5**（载体）＋**票 196**（状态词表，已 `-done`） | 并格；**不另立票** | 145 同上；196 已结案 | 行 4 的 chip 已把「运行中／成功／失败／被拒 ＋ L0/L1/L2」定义完，缺的是把 sink 事件搬进载体 |
| K2（可选）独立时间线 | **须 owner 一句话** | 若做＝改 `PLAN.md` 十四行表（新增一行） | — | 加行＝改契约文字，按铁律只摆 owner 批，本程不替裁 |
| K3 主体 | **票 145 行 6／7** | 不并（已在） | — | `pending` 已是快照第一维；145 表别重复列一行（防双记账） |
| K3 序号／停止 | **票 167 AC#3／AC#4** | 并格 | 1 勾／6 未勾 | `LiveApproval.Position` 已有真源（`pending_read.go:58/:127`），停止维是 AC#4 原文 |
| K4 文件树 | **票 190** | 已闭 | **`-done`**，6/6 | `TreeView`＋`files`＋生产装配都在；`tree_reader_190.go` 写调用 0 |
| K5 插件（读） | **票 407**（2026-11-05 新建，来路＝本普查） | 立票 | 无 AC 框；**Blocked by 前端 face 归外部 agent** | 全池唯一认领者；库里有读面、生产零装配、栏位规格真空 |
| K5 开关（写） | **无票；不进本票交付** | 单列待人拍板 | — | 写配置＋改档位＝`AGENTS §1.2` 那条禁令的邻域，见 §5 M3 |
| K6 内容审查 | **票 189** | 已闭 | **`-done`**，5 勾／1 未勾（AC#1 设计核） | `Uncommitted`／`Diff` 在同一枚 `GitView` 里加维，没另起实现 |
| K6 状态显示 | **票 181** | 已闭 | **`-done`**，7/7 | 分支／工作树／repoRoot／switchBlocked 全在 `git.go` |
| K7 工具地基 | **票 163 AC#1b** | 只登记依赖，不重开 | 1 勾／5 未勾 | 该票就是「`shell.session` 未注册」的账；注册循环 5 支无 Shell 是旁证 |
| K7 面板载体 | **票 191（甲支）** | 并格到该票 | 1 勾／6 未勾（AC#0 已裁甲） | 191 整张票就是「面板终端这一条执行面怎么裁」 |
| K8 子代理 | **票 197**（状态/流页）＋**票 188 AC#1**（已裁乙） | 并格 | 197：2 勾／5 未勾；188：1 勾／5 未勾 | `TaskRowView` 三件套是 197 落的 |
| K9 后台任务 | **票 188 AC#2**（状态维）＋**票 164**（已闭）＋**票 196**（已闭） | 并格 | 188：1 勾／5 未勾 | 名册在、状态维在 `TaskRowView.Status`；剩的是 188 AC#2 的口径更新（§9-2） |
| K10 浏览器（面板） | **票 192**（面板浏览器／产物预览那一堆） | 不立新票 | 1 勾／5 未勾（AC#0 已定案） | 192 的 AC#1–#5 写的就是「预览／注入面／只读边界／面板载体」，且排程明写与 191 串行 |
| K10 浏览器（能力/插件） | **票 168** | 不立新票 | 0 勾／7 未勾 | 168 管「逛网页」那门能力与插件门控；它自己逐字写「不做给人看的浏览器界面（界面腿 owner 委托）」⇒ 面板那一堆不是它 |

**⚠ 一处归属重叠，具名报回不替裁**：票 192 的标题与 §现场写「浏览器／产物预览**没票认领**」，而票 168 早已认领同一门能力（能力侧）；票 168 又把「给人看的浏览器界面」排除在自己写面外。⇒ 口径应为「**K10 面板堆＝192，K10 能力侧＝168**」，但两票的票面各自都把对方那一半说成「没票」，**这一句须编排者收敛**（本程不改票面）。

**⚠ 另一处归属是今天才变的**：K5 从「全池无票」变成「407 认领」。407 自陈 Blocked by 前端 face ⇒ **本编队的可写面只剩「Go 侧 bridge 读口若缺另立子格」**；本程没读 `internal/plugin/**` 的实现细节去替 407 出方案（那是它的活）。

### 3.1 AC#3 那个特别问：面板的「宿主侧读面」要不要与模型侧工具（`fs.list`／`shell.exec`）**共享实现**？

**先纠正 r1 的一条读数**（这是本程最硬的一处订正）：

| 事实 | 今天的尺 | 读数 |
|---|---|---|
| 面板生产代码 import 模型侧工具包 | `go list -f "{{join .Imports \" \"}}" ./internal/panel \| tr " " "\n" \| grep -c "internal/tools"` | **0**（权威尺；`go list` 只报真 import） |
| 同上，正控（同包确有一个 internal 包被 import） | 同一把尺 `grep -c "internal/risk"` | **1**（尺能响） |
| 面板文件里「internal/tools」出现次数 | `grep -rn "internal/tools" internal/panel --include=*.go \| grep -v _test \| wc -l` | **17 行**，逐行过目**全是注释**（`approval.go:7`、`git.go:24`、`pump.go:531` …） |

⇒ r1 §4.1 把「17 行」当成了「包在」的正控，**这是把注释读数当依赖读数**：面板**在生产代码里根本不依赖 `internal/tools`**，唯一的 import 在 `internal/panel/workspace_handlers_test.go:21`（测试）。模型侧工具的类型与判级是**经 `cmd/wisp` 装配点以值/字符串形式递进来**的（`run.go:927-960`），不是包依赖。

**两败（摊清，本腿不替编排者选，但给出读数结论）**：

- **共享**：`fs.list` 带 D34 判级（L0／L2 越界）与 C26 allowlist 语义（`PLAN.md:2540` 那一行）。面板一复用，**「面板显示」就走进「模型被允许读什么」的判定**；而 `AGENTS §1.2` 逐字禁的是「由面板侧来源的 L2『允许』」——共享会把禁令射程从「批准」扩到「读」，**需要人工批准**。
- **不共享**：两份枚举会漂（「面板看到的树」与「模型被允许读的树」不是一棵）。票 189 AC#2 用「同一文件族不另起实现」防的是 **git 一族内部**，不是跨「面板 vs 工具」。

**今天的既成事实（读数，不是建议）**：三枚已落地的宿主读面**各自独立**——181 用文件读（`git.go` 全文无 `exec.Command`）、189/190 各租一次 `git`（`git_readonly_exec.go` 固定查询字面，写动词 0）、190 的文件树不经过 `fs.list`（面板侧 `fs.list` 命中 0）。⇒ 现状＝**不共享实现、只共享路径/门控语义**。K7 的 191 甲支是「复用既有工具腿＋门控＋审计」，那是**共享门控、不共享显示通道**的第三条路，与 181/189/190 同向。
**本腿的建议（判的，不是量的）**：今后所有宿主读面按这条先例办（只读、独立实现、路径决策走 `risk.PathResolver`）；但这句按票面 AC#3 是**交编排者/owner 裁**的，本腿不把它当定案。

---

## 4. AC#4 与票 181 的分工（写死，两票今天都已结案）

| 维度 | 归属 | 今天在盘上的读数 | 框态 |
|---|---|---|---|
| git **状态显示**（是不是 git 树／当前分支／detached／本地分支名册／主工作树） | **票 181** | `internal/panel/git.go:144-171`（`Kind/Reason/Branch/DetachedSha/IsDetached/RepoRoot/CurrentWorktree/Worktrees/Branches/SwitchBlocked`） | **`-done`**，7/7 |
| git **内容审查**（未提交枚数／逐文件 diff／改了哪几枚文件） | **票 189** | 同一枚 `GitView` 的 `Uncommitted`（`:178`）＋`Diff`（`:181`）；`GitDiffView.files/filesTotal/filesIncluded/byteCap` 与 `GitDiffFile` 10 枚 | **`-done`**，5 勾／1 未勾（AC#1 设计核） |

- **两票都已 `-done`** ⇒ 182 AC#4 担心的「两票互相归口到对方已勾的格充数」今天**无处可归**：本程没有做、也没有建议做这种归口。
- **一处写死的边界**（从 189 的注释读出来的）：`git.go` 与 `git_readonly_exec.go` 是同一族两个文件，新增的只有 `composer.git` 底下的两个**嵌套子段**（`uncommitted`／`diff`），**没有新增顶层键**——这是 189 为避开 `Q-51`（前端键集契约）刻意选的位置。

---

## 5. AC#5 雷区清单：面板发起的批准/写动作

**现状合规锚（现读）**：面板入向方法名册 **7 枚**——`internal/panel/bridge.go:42-45` 四枚 `panel.*`（`mode.request`／`workspace.request`／`attachment.add`／`message.send`）＋`:66-67` `config.get`／`config.set` ＋`:75` `config.models.list`；`knownComposerMethod` 逐个点名同样 7 枚，**里面没有 decide／allow／commit／push／revert／stop／remove／apply／toggle／terminal**（`grep -nE "\"(decide|allow|approve|commit|push|revert|stop|remove|apply|toggle|terminal)" internal/panel/bridge.go`＝**1 行，且是 `:65` 注释里逐字引禁令**）。
`config.set` 另有一道闸：`lockedFieldFamilies`（`config_handlers.go:109`）逐字含 `plugins.`／`risk.`／`fs.`／`net.` 等族 ⇒ **插件启用开关走不到设置页的写入腿**。

| # | 雷 | demo 控件（现读） | Go 侧今天的门 | 定性 |
|---|---|---|---|---|
| M1 | **附件移除** | `rightbar.js:53`／`:74` `data-attach-del`（2 处） | 入向只有 `panel.attachment.add`，**无 remove** | **现存写形状，待拍板** |
| M2 | **diff 接受／应用** | `rb-review.js:96` `#rb-apply-btn`「Apply 3 changes」（`:115-119` 只 toast） | 入向无 apply／accept；189 已钉「不许出现动作腿」 | **现存写形状，待拍板** |
| M3 | **插件启用开关** | `rb-plugins.js:65`／`:75` `data-plg-toggle`（2 处，条目带 `on:`／`level:`） | 入向无 plugin 方法；`config.set` 把 `plugins.` 列为锁定族 | **现存写形状，待拍板**（写配置＋改档位，最硬的一枚；读面归 407） |
| M4 | **面板里发起执行** | `rb-terminal.js:26`／`:77` `data-terminal-add`（2 处）＋命令输入行（`:86-` 只本地回显） | **地基都不存在**：注册循环 5 支无 Shell | **现存写形状，待拍板**；191 甲支是它的出口 |
| M5 | **面板「允许」** | 今日 demo **无此控件**：approval 条目只 `showScreen('approval')`（`rightbar.js:143`／`:154`／`:160`）；`grep -nE "允许\|拒绝\|批准\|approve\|allow" rightbar.js`＝**2 命中，都是散文**（`:93` 时间线条目、`:136` 满额告警） | 入向名册无 `decide`／`allow`；`bridge.go:63-65` 注释逐字引 `AGENTS.md §1.2` | **潜在雷（今日无控件）**：任何新增即逐字撞禁令＋`Q-49` 一族 |
| — | ~~文件树删除／移动~~ | `grep -nE "删除\|移动\|重命名\|trash\|data-del\|data-move" rb-files.js`＝**0** | 190 钉了只读；`tree_reader_190.go` 写调用 0 | **今日已消失** |
| — | ~~面板 git 提交／推送~~ | `grep -nE "接受\|回滚\|提交\|推送\|commit\|push" rb-review.js`＝**0** | 入向名册无 commit／push；`git_readonly_exec.go` 写动词字面 0 | **今日已消失** |

**现数 = 4 枚现存写形状（M1–M4）＋1 枚潜在（M5）**；09-28 记过的两枚（文件删除/移动、git 提交/推送）今天在 demo 与 Go 两侧都读不到控件/入向。
⚠ M1–M4 **一律单列待人拍板**：`AGENTS.md §1.2` 逐字禁「由面板侧来源的 L2『允许』」，`Q-49` 一族已钉过五枚门；
**本程没有提「顺手把批准接上」**，也没有把它们并进只读交付。

---

## 6. 返图：145／187／188／192／195 这五张"扣着"的票，各自欠本题哪几堆

> 这一节是派单里「return-map」那句的落点：五票都还没结，本表说清**谁欠哪一堆、欠在哪一格**，防止下一轮又拿局部截图重新问 owner。

| 票 | 框态（现读） | 在 10 堆里**持有** | 明确**不持有**（不许归口过来充数） | 今天卡在哪（具名） |
|---|---|---|---|---|
| **145** 快照扩字段 | 4 勾／3 未勾（AC#2／AC#2b／AC#6） | **K1-D**（行 12 成本）、**K2**（行 4／5 扩载体）、**K3 主体**（行 6／7，`pending` 已在，别重复列） | K4（190 已闭）、K6（189 已闭）、K7（163/191）、K8/K9（197/188） | AC#2 的量出答案是「落地集＝空」（145 自己的收表），AC#6 反向判据要等一枚可判归宿；票面自陈挂在 `Q-51` 与 R1/R2 |
| **187** 模型/思考档位可改 | 0 勾／9 未勾 | **无**（它在底部组合栏，不在右上竖条；本程**没有**把 K1 的任何一行算给 187） | 本题全部 10 堆 | 需要先答「改了哪一档生效」与写通道（C17 新方法＝契约变更）；其**显示半**已由 `ComposerState.currentModel/modelKnown`（`composer.go:290-291`）落，属 145 AC#2b 的一半 |
| **188** 子代理/后台任务状态 | 1 勾／5 未勾（AC#2–#6） | **K8**、**K9** | K1-C（167）、K10（192） | AC#2 前提写「给 `Record` 补状态维」，而 `type Record struct` 在本树 **0 命中**（正控 `type TaskRoster struct`＝**1**，`internal/tools/task.go:176`）⇒ **该格口径过期**；状态三件套已由 197 落在 `TaskRowView`（§9-2） |
| **192** 浏览器/产物预览 | 1 勾／5 未勾（AC#1–#5） | **K10 面板那一堆**（浏览器标签＋产物预览） | K7 终端（191）、浏览器**能力侧**（168） | AC#0 已定案；AC#1–#5 未勾；排程明写：`33-r1` → 191 终端落地腿 → 192 预览 → 192 联网（三枚同包串行） |
| **195** `config.set` 分节写地基 | 0 勾／7 未勾，**`blocked`** | **无只读堆**（它是写地基） | 本题全部只读堆 | 票面 `blocked`：排在票 33 片 B 交件之后；`33` 今天**没有 `-done`**（`ls .scratch/wisp/issues/ \| grep -E "^33-"`＝`33-panel-host-c27.md`）⇒ 前置仍在 |

**这张返图的结论（一句话）**：五票里**只有 145／188／192 三票持有本题的堆**；187 与 195 在本题里**一堆都不持有**（前者在底部组合栏，后者是写地基）——下一位若拿「任务监控那一栏」去催 187／195，那是拿错票。

---

## 7. RED-before / GREEN-after：本件相对 r1 与票面的六处订正（每处都是「旧断言在今天的树上响不了」）

> 只读票没有产码，RED 只能是**旧断言被尺打红**。下表每行：旧断言（出处）→ 今天的尺 → 读数 → 判。

| # | 旧断言（出处） | 今天的尺 | 读数 | 判 |
|---|---|---|---|---|
| 1 | 票面现量#1：「票池里**没有一张票认领它**」 | `grep -rlniE "内置终端\|文件树\|任务监控" .scratch/wisp/issues/*.md \| wc -l` | **8**（77／182／186／188／189／190／191／196；与 r1 同数同批，无第 9 枚） | **仍过期**（立票链早接住） |
| 2 | r1 §3：K5「**无票认领**」 | `ls .scratch/wisp/issues/ \| grep "^407-"` | `407-k5-plugin-panel-has-read-surface-no-ui-no-spec-unowned.md` **存在**，2026-11-05 由编排者按 owner「立新票」建 | **已过期**（K5 归属改为 407；407 自己被前端 face 扣住） |
| 3 | r1 §3：K10 浏览器 → **票 168** | 读 192 票面标题与 §现场、168 票面「本票**不**解决」段 | 192 逐字「浏览器／产物预览那一堆没票认领」；168 逐字「不做给人看的浏览器界面（界面腿 owner 委托）」 | **过期且重叠**（面板堆＝192，能力侧＝168；两票须编排者收敛，见 §3 表脚注） |
| 4 | r1 §4.1：面板「从 `tools` 包拿的只有类型」＋正控「`internal/tools` 17 行＝包在」 | `go list` 权威尺 | **生产 import 0**；17 行**全是注释**；唯一 import 在 `workspace_handlers_test.go:21` | **读数错**（依赖结论反而更干净：面板与工具包无代码共享） |
| 5 | 票面现量#4 与 r1 §10：「`git diff` 源不在」 | `grep -rn "git diff\|diff --numstat" --include=*.go internal/ cmd/ \| grep -v _test.go \| wc -l` | **15 行**（全部在 `internal/panel/git_readonly_exec.go`） | **已过期**（189 落了内容审查） |
| 6 | 168 的 Status：「前置＝50＋51 **都没做完**，插件地基不存在」 | `ls .scratch/wisp/issues/ \| grep -E "^50-\|^51-"` | `50-tier1-manifest-plugins-done.md`、`51-tier2-goja-done.md` **两枚都带 `-done`** | **前置已满足**⇒ 168 今天真正卡的是它自己要动的 `C17`/`C24` 白名单**未定案（人工批准）**，不是地基 |

**GREEN 的那一半**（今天跑出来的、与旧断言相反的正读数）：`GitView` 12 枚 json 含 `uncommitted`/`diff`；`TaskRowView` 12 枚含 `status`/`statusKnown`/`statusReason`；`TreeView` 7 枚＋`files` 键＋`run.go:932` 装配；`ComposerState` 12 枚；`Snapshot` 6 键；`PumpSources` 11 个 reader（无 Cost／无占用）。
**没有一处是「放宽断言换绿」**：本程不写码、不改测试、不翻框；所有 0 命中都带同包正控。

---

## 8. 门禁读数（本程现跑，命令与 rc 原样）

| 门 | 命令 | 读数 |
|---|---|---|
| build | `export PATH="$PWD/build:$PWD/third_party/sherpa-onnx:$PATH"; unset GOTMPDIR TMPDIR TEMP TMP; go build ./cmd/... ./internal/... ./tools/...` | **rc=0** |
| test（主题包） | `go test ./internal/panel/ -count=1 -timeout 420s` | **rc=1**，两枚红，**皆具名且皆非本程造成**（`git status --porcelain -- internal/panel`＝空）：① `TestComposerContractTypesMatchFrontend` 报「Go ComposerState emits [files] that interface ComposerState does not declare」＝K4 的**前端声明欠账**（`Q-51` 一族，190 的 `files` 键）；② `TestC21DesignTokensFourWayAgree` 报 `c21-native-tokens.md` 与 `internal/ball/tokens.go` 6 枚 `Palette.Glass*` 不同步（票 77 的 AC#2 段与 145 的收表都已记过这把红） |
| test（邻包，为 K9 的源） | `go test ./internal/tools/ -count=1 -timeout 420s` | **rc=0**（`ok github.com/CarlosShao/wisp/internal/tools 24.845s`） |
| gofumpt（产码 scope） | `/d/work/base/gopath/bin/gofumpt.exe -l cmd internal tools` | **空，rc=0** |
| gofumpt（CI 全仓形） | `gofumpt -l . tools/d22scan tools/mockllm` | **rc=2**：52 枚，全部住在 `.scratch/**`（50 枚）＋`build/probe61r2/main.go`（1 枚）等**他票探针/变异件**，`cmd/**`／`internal/**` 命中 **0**（票 276 在册：这批红来自 14 张别家工单的探针件，本程未动、也不该顺手格式化） |
| d22scan | `sh scripts/d22scan.sh` | **rc=0 clean**：`examined 421 production Go files`；scope `bans#1-5 internal/=363 cmd/=58`、`ban#6 frontend/=116`、`ban#7 internal/tools/=48`、`ban#8 design/=39 frontend/=116 internal/=905 cmd/=212`；自测 `PASS=35 FAIL=0 SKIP=0` |

⚠ 全仓 gofumpt 的红**不是**「本腿放宽门禁」：同一把尺在产码 scope（`cmd internal tools`）上是空的；全仓红逐枚落在 `.scratch/**`／`build/**`，且票 276（open）就是它自己的账。

---

## 9. 缺陷与未判项（具名，不许当完整表读）

1. **〔越界自报〕本程读了一次 `frontend/**`**：为定位 `TestComposerContractTypesMatchFrontend` 的红因，对 `frontend/src/lib/panel.ts` 跑过 **1 次 `grep -n "files"`**。票 182 AC#7 逐字「`frontend/**` 不读不写不引不转述」⇒ **这条我越了**。补正：该红因**不需要那条 grep**——测试自己的报错行已经逐字给出「Go ComposerState emits [files] that interface ComposerState does not declare」。本程没有写、没有引、没有转述该文件内容，本件只保留测试报错那一行。
2. **票 188 AC#2 的前提与代码对不上（具名报回，不替裁）**：该格逐字写「给 `Record` 补状态维」，而 `grep -rnE "type Record struct\|type Job struct\|type BackgroundJob" internal --include=*.go \| grep -v _test.go`＝**0**；正控 `type TaskRoster struct`＝**1**（`internal/tools/task.go:176`）。状态维真身在 `TaskRowView.Status/StatusKnown/StatusReason`（197 落的）⇒ 该格要么改口径、要么改归属，**须编排者裁**；本程不改票。
3. **K5 的「有源」只到库层**：`memory.ListPluginStates`（`dao_misc.go:243`）只有定义＋两处测试调用，**非测试调用者 0**；`grep -rn "plugin/manifestslot" internal cmd --include=*.go \| grep -v _test.go`＝**1 行注释**（`internal/plugin/command.go:16`），正控 `internal/plugin"` 非测试 import＝**9 行** ⇒ `manifestslot` **生产 import 0**。所以我记的是「**库里有读面、生产零装配**」，不是「运行时能取到真值」。
4. **K4／K6 的读面都租了 `git` 可执行文件**（`tree_reader_190.go:27`、`git_readonly_exec.go:349`）：这是**读数**，不是我对「应不应该起进程」的判断；189 AC#1 那一格（设计核）框未勾，谁要评议该去读 `189-uncommitted-count-design-a1.md`。
5. **K1-C 的分子/分母口径本程没量完**：只量到「`internal/agent/budgets.go` 有 `TokenBudget`/`ContextWindow` 两枚分母候选、`Result.Usage` 是分子候选、面板零字段」。demo 的「42%」到底除的是什么，**今天没有任何一节写死** ⇒ 记「口径未定」，不记「无源」。
6. **`PLAN.md` 行号漂移**（引用者必读）：票 145 与 145 证据件引的 `:3473-3488`／`:3481-3494` 今天**都不是十四行表**；现量＝表头 `:3497`、十四行 `:3499-3512`、**行 12「成本」＝`:3510`**；D34 的 `fs.list`＝`:2540`、`shell.session`＝`:2575`。
7. **`Snapshot` 的 6 键与「线上 4 键」不是一回事**：`instructions`／`tasks` 是 `*T,omitempty`，没有 reader 的泵**不送这两键**（`pump_test.go:123/:291` 的两枚字节钉钉的就是这个形状）⇒ 引用键集时要说清是「声明 6」还是「某一次装配实际送出几枚」。本程记的是**声明 6**。
8. **`182-a1`/`182-c1` 转引的三处本程未复跑**：K9 的 `rightbar.js:90`、K10 的两处 `issues/77` 出处（本程**复跑了**）、K2 的细维出处。其余转引处一律标「转记」。

---

## 10. 一句话交件

**10 堆里：整堆已落 4 堆（K4 文件树／K6 审阅／K8 子代理／K9 后台任务）；部分已落 2 堆（K1 的路径与附件、K3 的 `pending` 主体）；
只差票 145 那一集载体 2 处（K1-D token／K2 活动）；只差票 167 三个维 3 处（K1-C 占用／K3 序号／K3 停止）；
无源但票已在手 1 堆（K7 终端，163 地基＋191 甲载体）；面板堆认领已改 2 处（K5→407，K10→192＋168 能力侧）。**
雷区 **4 枚现存（附件移除／diff 应用／插件开关／面板执行）＋1 枚潜在（面板允许）**，一律单列待人拍板。
五票返图：**145 持 K1-D／K2／K3 主体；188 持 K8／K9；192 持 K10；187 与 195 在本题一堆都不持有。**
本程零产码、零翻勾、`internal/**`／`cmd/**`／`docs/PLAN.md`／`docs/specs/**` 零字节；写面只有本文件与票 182 的一行 Progress log。
