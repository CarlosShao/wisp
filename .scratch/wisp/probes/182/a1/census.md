# 182-a1（本轮·票 182）—「任务监控」右栏逐堆 Go 侧有无源：三档定性＋归属指认＋与票 181 分工＋雷区

- 编队：`182-a1`（只读普查，零产码）。派单＝编排者任务书（无 dispatches 文件，射程以任务书为准）。
- 起手锚：`e16e303773f59be81d2f02eee3cc638870aa1630`（`dev`），起手时刻 `2026-10-03T10:33:41+08:00`。
- 写面：只有本件 `.scratch/wisp/probes/182/a1/census.md`。`internal/**`、`cmd/**`、`docs/**`、`frontend/**`、`design/**` **零字节改动**；票面 AC 框**一枚没勾**；不 push。
- 起手闸门读数（带时刻）：`2026-10-03T10:37:02+08:00`
  - `git status --porcelain -- internal/ cmd/` ⇒ **2 枚**（`M cmd/wisp/config_reload.go`、`?? cmd/wisp/restart_tier_keys_255r2_test.go`）＝**在飞的 `255-r2` 写腿的**，不是本程的。本程自己名册差集＝空。
  - 全树 `git status --porcelain` ⇒ **473 枚**（` M`16／` D`16／`??`441）——同一时刻的数，含 `??` 未跟踪，见 §5。
  - `git rev-list --left-right --count cnb/dev...HEAD` ⇒ `0  115`＝**未推 115 枚**（只对本程有效，别家写腿一提交就漂；口径＝整棵树，不是"我这轮"）。
- ⛔ 本程**没跑任何 Go 命令**（`go test`／`build`／`vet`／`list`／`wisp` 一枚没跑）——三枚写腿在飞。凡"要跑才知道"的格子一律进 §6 量不到清单，不猜。
- ⛔ `frontend/**` 与 `design/**` **不读不引**（任务书 §2 两层禁令）。**后果写清**：本栏堆数的上一轮主尺是 `design/doubao/demo/` 的注册表，本程不能复跑 ⇒ 本程的堆数是**仓内文本＋Go 载体两向**并集口径，与上一轮"demo 注册表口径"**不同尺**（见 §7 推翻清单 P6）。

---

## 1. AC#1 原文：先给这一栏下定义并数出堆数

> 票面要求拆成可点名的堆、**枚数只许现量不许目测**、逐堆写"我在谁那儿看到的这一堆"，并明写 owner 那张截图是**外部产品界面、不是本仓规格**。

### 1.1 本程用的尺（三把，全部允许根，逐把可复跑）

| 尺 | 命令（逐字） | 读数 | 用它干什么 |
|---|---|---|---|
| R-A 认领尺 | `grep -rliE "任务监控\|右栏\|那一栏" .scratch/wisp/issues` | **18 枚票**命中（154/164/181/182/186/188/189/191/192/196/219/246/36/37/38/39/40/77） | 数"有多少张票在谈这一栏"，**不是**堆数 |
| R-B 载体尺 | 现读 `internal/panel/composer.go:57-92`（`Snapshot`）＋`internal/panel/pump.go:138-207`（`PumpSources`） | 顶层 **6 枚键**／泵读口 **10 枚** | 数"这一栏今天真能读到什么" |
| R-C 名册尺 | `grep -rniE "文件树\|内置终端\|任务监控\|子代理\|产物预览" docs/specs/ docs/PLAN.md` | `docs/specs/**` **1 命中**（`SPEC-01:28`＝仓库骨架那棵树，**不是界面**）；`docs/PLAN.md` **4 命中**（`:7`／`:1123`／`:1127`／`:1544`，全是"不做子代理"那一族） | 判"文档有没有要求过这一堆" |

⚠ 上一轮（`docs/evidence/s1/182-rail-stacks-census-a1.md` §2）的主尺 `grep -rn --include=*.js "RbPanels.register(" design/doubao/demo/` ⇒ 本程**未复跑、不可复跑**（禁读）。那一轮的 7＋3＝10 枚在本程只作为**被转记的既有读数**出现，逐处标注。

### 1.2 堆数现量：**12 枚**

把前两程（`182-c1` 的 S1–S7 ＋ `182-a1` 的 K1–K10）按**内容去重**（不是按编号对齐），再用今天的 Go 载体尺过一遍，可点名＝**12 枚**。去重过程逐枚写：`182-c1` 的 S1「子代理／后台任务」在本程拆回两枚（它有两只不同的手：`PublishSubagent` 与 `MarkRoot`）；`182-c1` 的 S7「成本 token」在上一轮被并进了 K1 的"上下文预算"那一行，本程**拆开数**（两者源不同、缺的不同，见 §2 行 4／行 12）。

| # | 堆（可点名） | 我在谁那儿看到的这一堆（逐字出处） | 是不是本仓规格 |
|---|---|---|---|
| A1 | 子代理状态（每一枚的 kind／父／状态／自己的工作页） | 票 188 标题逐字 "subagent and background task stacks have no go-side source"；票 197 owner 原话逐字「子代理必须看到状态，而且点击某个子代理，能看到它们各自的流式工作页面」 | **不是规格**，是 owner 口头放权（`A373`／`Q-71`）＋票 197 接单 |
| A2 | 后台任务名册（这一程跑了哪几枚任务） | 票 188 `AC#2`；票 164 标题逐字 "Background jobs cannot be read — fix the roster then add the output leg" | 同上 |
| A3 | 环境信息（工作目录／分支／工作树／**未提交枚数**／提交或推送） | 票 182 `AC#1` 目测清单逐字「环境信息（未提交枚数／本地路径／分支／提交或推送）」；票 181 只读那半已落地 | **半规格**：分支／工作树＝票 181；未提交／未推＝票面目测清单，⚠ 截图是**外部产品（qoder）界面，不是本仓规格** |
| A4 | 上下文那一行（附件清单／接受类型／上下文预算） | 票 92（附件四枚键）；票 214 标题逐字 "attachments and context wired"；`182-a1` 转记 K1 三行 | 票面级有、`docs/specs/**` 无 |
| A5 | 活动／时间线（工具调用行＋状态＋耗时） | `182-a1` 转记 K2；票 145 那张十四行表（`PLAN.md:3484` 工具调用行） | **规格有**，但落在票 145 那一族、不在本栏 |
| A6 | 审批队列（哪几枚在等／队列深度） | 票 145 `AC#1` 行 6/7（`PLAN.md:3486`／`:3487`）；票 146／220 已落；`Snapshot.Pending` 在册 | **规格有**（且这一堆本来就在 145 射程里） |
| A7 | 工作区文件树 | 票 190 标题逐字 "the workspace files stack file tree has no host-side source"；票 77 `:384` 逐字含"审查/终端/浏览器三面板、任务监控弹窗" | `docs/specs/**` **零要求**（R-C 唯一命中是仓库骨架那棵树） |
| A8 | 插件／技能清单（列出装了哪些、开关、生效没有） | `182-a1` 转记 K5；票 215（`215a` 插件／`215b` 技能） | `docs/specs/**` 零要求；票 215 接单 |
| A9 | 代码审查（逐文件 diff／改了哪几枚文件） | 票 189 标题逐字 "uncommitted count and per-file diff but go-side has zero source"；票 77 `:364` 逐字"右栏三标签定案：审查/终端/浏览器" | 票面级有、`docs` 无 |
| A10 | 内置终端 | 票 191 标题逐字 "the panel embedded terminal puts an execution surface on the panel channel"；`PLAN.md:2569` **D34 表里有 `shell.session` 这一行** | **规格要求，且是契约级**（D34 行在册） |
| A11 | 浏览器／产物预览 | 票 192 标题逐字 "the browser and artifact preview stack is unclaimed"；`PLAN.md:2546-2548` D34 三行 `web.search`／`web.fetch`／`web.open` | **拆两半**：浏览网页＝规格要求且压在 `AGENTS §2` 未定案项上；产物预览＝`docs` 无 |
| A12 | 成本与 token（`IN`／`OUT`／`CACHED`／金额） | `PLAN.md:3492`＝§17.5 那张表**第 12 行**逐字「**成本**（C23）｜纯文字…`12,480 tok · ¥0.31`…标注 `IN`/`OUT`/`CACHED`」；票 44（`costmeter-c23`） | **规格有**（C23＋那一行就是票 145 的行 12） |

⇒ **本程结论：这一栏可点名的堆＝12 枚**（前两程现量的 7 与 10 都不是错，是**两把不同的尺**：`182-c1` 数"owner 目测清单＋仓内文本"、`182-a1` 数"demo 注册表＋有名无位"，本程数"去重后的可点名数据堆"）。三枚口径都记在案，谁复跑都知道自己用的是哪把尺。

---

## 2. AC#2 原文：逐堆三档定性（有源／无源但规格要求／规格真空）

> 硬要求：每档附尺与命中数；0 命中先拿一枚已知存在的正控打过那把尺。
> 判"接上了没有"一律问三处——**谁产出／谁投递／谁落盘**——不数"显式调用者枚数"。

**汇总（算术写清，别把半堆混进整堆）：12 堆＝**
- **有源（缺载体或缺投递）＝5 整堆＋1 半堆**：A1、A2、A5、A6、A12 ＋ A3 的分支／工作树那一半。
- **无源但规格要求＝2 整堆＋1 半堆**：A10、A9 ＋ A11 的浏览网页那一半。⚠ A9 的"规格"只到**票面级**（`docs/specs/**` 对 diff 零要求），把它放这一档是**从宽**，从严则落"真空"——本程两话都写，见 §8 表内注。
- **规格真空＝2 整堆＋1 半堆**：A7、A8 ＋ A11 的产物预览那一半。
- **半堆加总校验**：只有 **A3** 与 **A11** 两枚整堆被拆档 ⇒ 不拆的整堆＝5（有源）＋2（无源有规格）＋2（真空）＋1（A4）＝**10 枚**，加回 A3、A11 两枚＝**12**，无重复计数。
- **另记一档＝1 堆**：A4 上下文那一行（**载体在、投递被硬写 `nil`、预算子维无源**）⇒ 它不属于上面三档任何一档，第四种形状＝"接了半根线"。

⚠ **不许把"有源"读成"接上了"**：A1／A2／A6 的"接上"三处都点到了名；A12 那一族本程只证到**没有读口**，至于"补一枚 `Cost func() ...` 读口会不会红那三枚在册的字节级键集尺"——**量不到**（§6 M1）。

### 2.1 逐堆判定

| # | 堆 | 三档 | 谁产出 | 谁投递 | 谁落盘 | 尺与命中数 |
|---|---|---|---|---|---|---|
| A1 | 子代理状态 | **有源＋已有载体** | `internal/tools/subagent_197.go:310`／`:345` `PublishSubagent`（＋`:270 TryAcquireSubagentSlot`）三处都有生产调用者 | `cmd/wisp/run.go:724 Tasks: rt.taskRosterState` → `cmd/wisp/panel_pump.go:146` | `Snapshot.Tasks`（`internal/panel/composer.go:91`）→ `TaskRowView`（`internal/panel/subagent_roster_197.go:106`） | 目标尺 `grep -rln "ubagent" --include=*.go internal cmd \| grep -v _test.go`＝**9 枚文件**（⚠ **直接推翻票 197 现量表"0 命中"**，那一行是 197 **落地前**的数，见 §7 P1） |
| A2 | 后台任务名册 | **有源＋已有载体** | `cmd/wisp/run.go:1105 MarkRoot` ＋ `internal/tools/subagent_197.go:344 RunAsync`；`Record` 写点 `internal/tools/task.go:215`（经 `:343`／`subagent_197.go:424` 两路进真值） | 同上（`taskRosterState` 一次性交 Rows＋`InFlightSlots`＋`PoolCap`） | `TaskRosterSection`（`subagent_roster_197.go:150`，含 `inFlightSlots`／`poolCap`／`droppedStreamKeys`） | 状态维：`internal/tools/task.go:158 StateAnswer()` **存在且校验 D43 的 20 枚名字**（`:164` 逐字"不是 D43 状态机表里的 20 个名字之一（票 188 AC#2：不许自造）"）⇒ 上一轮"名册在、状态维没有"**已过期**（§7 P2） |
| A3 | 环境信息 | **一枚堆里三档并存**（不许整堆算一档） | 分支／工作树：`internal/panel/git.go:172 ReadGit`（零外部进程，正控尺 `grep -rn "ReadGit"` 非测试 **10 命中**）→ `ComposerState.Git`（`composer.go:246`）＝**有源＋已有载体** | 未提交枚数：**无任何产出者** | 未推枚数：**无任何产出者** | 未提交尺 `grep -rniE "git status\|porcelain\|numstat\|git diff" --include=*.go cmd internal tools \| grep -v _test.go`＝**0**，正控同形尺（`ReadGit`）＝10 ⇒ 尺能响。**载体侧同判**：`GitView`（`git.go:117-145`）**10 枚字段里没有一维**是"改动枚数／差异内容／领先落后" |
| A4 | 上下文那一行 | **有载体、投递断、部分无源** | 附件本体：`internal/panel/attachments.go:179 NewAttachmentBroker`——**生产零调用者**（尺 `grep -rn "NewAttachmentBroker" cmd internal tools docs scripts` 非测试＝**只有那一行定义**，其余全在 `attachments_test.go`；件 410 行，与票 214 标题"410 行、生产零调用者"逐字同值） | `internal/panel/pump.go:269 NewComposerState(mode, workspace, nil, maxAttachment)` ⇒ 第三参数**硬写 `nil`** | 键在（`ComposerState.Attachments`／`AttachmentMIMEs`／`MaxAttachmentB`／`AttachmentReason`），值永远是空 | 上下文预算（用了多少／总共多少）＝**无源**：`grep -rniE "occupancy\|ContextBudget\|BudgetUsed" --include=*.go internal/agent internal/panel cmd/wisp` 非测试＝**0**（正控：同根 `grep -rniE "\bModel\b"`＝40），归票 167 |
| A5 | 活动／时间线 | **有源缺载体（半堆）** | 文本流：`cmd/wisp/run.go:998 Sink: consoleSink{...stream: rt.stream...}`（`internal/agent/loop.go:1033` 逐事件 `Publish`）；逐行页地址 `TaskRowView.StreamKey` | `PumpSources.Results`（`pump.go:165`）已投文本 | 已落 `Snapshot.Results` ⇒ **流式文本这一半不缺** | 工具调用行（名称／参数／耗时／成功失败）**载体无**：`turn.ToolCalls` 只活在 `internal/agent/loop.go:459-608` 的 run 局部；`internal/agent/sink.go:79 RecordSink` **生产零调用者**（尺 `grep -rn "RecordSink\|agent.Sink" cmd internal` 非测试＝0，正控＝同尺去掉 `Record` 换 `Sink:` 命中 2 枚）⇒ **归票 145（十四行表），不归本栏另立票** |
| A6 | 审批队列 | **有源＋已有载体** | 队列真值 `PumpSources.Verdicts`（`pump.go:140`）＋`L1Windows`（票 220 AC#2 的枚举） | 同一枚泵 | `Snapshot.Pending`＋`TaskRowView.BlockedOnApproval`（`subagent_roster_197.go`，读同包那枚 `pending`，两节不可能互相打脸） | ⚠ **这一堆本来就在票 145 射程里**（行 6/7）⇒ 派单 §0 那句"整块不在射程"**不成立**（§7 P4） |
| A7 | 工作区文件树 | **无源（宿主侧）** | 宿主侧尺 `grep -rn "os.ReadDir\|filepath.WalkDir" --include=*.go internal/panel cmd/wisp` 非测试＝**3 命中**，逐枚看过：`git.go:22`（注释）、`git.go:455`（refs 逐级走）、`git.go:498`（`.git/worktrees` 枚举）⇒ **三枚全在 git 那一族里，没有一枚是"给面板列工作区"** | 无 | 无 | 模型侧**在册**：`internal/tools/fs.go:190 fs.list`＋`BuiltinFSEntries`（`fs.go:325`）＝两回事（授权语义不同，见 §3.2） |
| A8 | 插件／技能清单 | **规格真空＋无源** | `internal/plugin/` 整包只有 `disposal.go` 一枚非测试件，导出的全是 `DisposalScope`／`MemoryReleaser` 那一族 ⇒ **零"列出已装"读面** | 无 | 无 | 配置侧确有键（正控）：`grep -c -i plugin internal/config/schema.go`＝**18**、`grep -c -i skill`＝**0**（与票 215 那句"8 处"不同＝它数的是另一把尺 `grep -c Plugin`（大小写敏感）；本程 18 是大小写不敏感，两把尺都对、口径要写清）。面板侧尺 `grep -rniE plugin --include=*.go internal/panel cmd/wisp` 非测试＝**3 命中**，逐枚：`config_handlers.go:108`（**锁死名册**里的一行）＋`approval_always.go:6`／`config_reload.go:126`（注释）⇒ **零读面** |
| A9 | 代码审查 diff | **无源但规格要求（票面级；`docs` 级无）** | 同 A3 的未提交尺＝**0** | 无 | 无 | 唯一现成的 git 起进程点在**另一个 module**：`tools/d22scan/gitignore.go:373 exec.CommandContext`，用途是 `ls-files`（`:274`／`:278`），**不是 status／diff**；`tools/d22scan/go.mod` 逐字 `module github.com/CarlosShao/wisp/tools/d22scan` 且 `internal/**`＋`cmd/**` 非测试对它**零 import**（尺 `grep -rn "d22scan" --include=*.go internal cmd` 非测试＝0 枚真 import，命中全是注释） |
| A10 | 内置终端 | **无源但规格要求（契约级：D34 行在册）** | 尺 `grep -rnE "\"shell\\.(exec\|session)\"" --include=*.go cmd internal` 非测试＝**0**，正控（同形，换成 `"task.(spawn\|output\|cancel)"`）＝**6**；`ls internal/tools` 非测试**没有 `shell*.go`** | 无 | 无 | 规格逐字：`PLAN.md:2569`「**`shell.session`**｜在一条活着的会话里跑下一条命令（工作目录与环境延续）｜**L2**｜`shell`｜S3｜票 163（2026-09-27 owner 批准新增）」＋`:2568 shell.exec`（D14 默认禁用）。⚠ **比票面更黑**：票 182 现量第 4 条说它"未注册"，本程量到的是**连工具实现体都没有**（§7 P3） |
| A11 | 浏览器／产物预览 | **拆两半**：浏览网页＝无源但规格要求（且压在未定案项上）；产物预览＝规格真空 | 工具实现体尺 `grep -rn 'func (.*) Name() string { return "web' --include=*.go internal cmd`＝**0**；字面尺 `"web.` 非测试＝**3 命中**，逐枚看过全在 `internal/risk/provenance.go:88`／`:120`／`:137`＝**名字册与盖戳通道，不是实现**（这三枚的存在正是正控：说明 `web.search`／`web.fetch` 这两个名字**在风险面上已被登记、在能力面上还不存在**） | 无 | 无 | 规格：`PLAN.md:2546`（`web.search`，逐字"实现路径仍待确认（§15 第 3 条）"）／`:2547`（`web.fetch`）／`:2548`（`web.open`，L1）；`AGENTS §2` 逐字列「`web.search` 实现路径（抓结果页 vs API，S3 前）」＝**未定义即停项** |
| A12 | 成本与 token | **有源缺载体** | `internal/agent/cost.go:38 AddUsage` → `internal/agent/loop.go:435 delta := cost.AddUsage(...)` → `:436 guard.AddCost(delta)`（`guard.go:155`）→ `:437-442 publishUsage`（`:976`），值＝**run／turn 局部** | **泵侧零读口**：`PumpSources`（`pump.go:138-207`）10 枚读口逐枚名过——Verdicts／Mode／Workspace／Git／Model／Credential／Results／Instructions／Tasks／L1Windows，**没有 Cost／Usage 那一枚** | **快照里零键**：`Snapshot` 6 枚、`ComposerState` 11 枚逐枚名过，无成本维 | 尺 `grep -rniE "\bCost\b\|price\|micros" --include=*.go internal/panel cmd/wisp/panel_pump.go` 非测试＝**13 命中**，逐枚看过：全是 `model_price_in`／`model_price_out` 两枚**设置页价格卡字段**（`config_handlers.go:61-62`）与注释里的英文单词 ⇒ **不是运行成本**（正控＝同尺 `\bModel\b` 40 命中）。规格出处＝`PLAN.md:3492` 行 12 |

### 2.2 三档各几堆（交件要的那个数）

- **有源**（Go 已取得到真值；缺的是载体或投递）＝**6 堆**：A1、A2、A5、A6、A12 ＋ A3 的分支／工作树那一半。
  ⚠ 这一档**内部还要分两层**，别混：A1／A2／A6 **三处齐了**（产出＋投递＋落盘都有名）；A5／A12／A3-树是**有源缺载体或缺投递**（A4 附件更硬：载体在、`pump.go:269` 那一跳自己把值写死成 `nil`）。
- **无源但规格要求**＝**3 堆**：A9（票面级要求、非 `docs` 级；`docs/specs/**` 对 diff 零要求）、A10（**契约级＝D34 表 `PLAN.md:2569` 已有那一行**，缺的是实现）、A11 的浏览网页那一半（**契约级 D34 三行在册＋`AGENTS §2` 未定案项**）。
- **规格真空**（今天没有任何文档要求过它）＝**3 堆**：A7 文件树、A8 插件／技能清单、A11 的产物预览那一半。
- **混合／另记**＝**3 处**：A3 里"未提交枚数"与"未推枚数"两枚子项（`docs` 级零要求＝真空，且连源都没有）、A4 上下文那一行（键集在、投递被写死成 nil、预算那一维无源）。

⚠ **不许把"有源"读成"接上了"**：A1／A2／A6 的"接上"三处都点到了名；A12 那一族本程只证到**没有读口**，至于"补一枚 `Cost func() ...` 读口会不会红那三枚在册的字节级键集尺"——**量不到**（§6 M1）。

---

## 3. AC#3 原文：逐堆指认归属票（本票主要产出）

### 3.1 归属清单

| # | 堆 | 归属指认 | 凭据（票面 AC 现量） |
|---|---|---|---|
| A1 | 子代理状态 | **票 197（已接单、已落载体）＋票 188 `AC#2`（状态维）** | 197 Status＝**在派**（2 勾／5 未勾）；188 checked=1／unchecked=5 |
| A2 | 后台任务名册 | **票 188（名册＋状态）＋票 164（输出腿）＋票 196（两副状态词汇）** | 164 Status＝**立而不派**（3 勾／3 未勾）；196 在册 |
| A3 | 环境信息·分支／工作树 | **票 181（已落地）**——本程现读到，不是纸面分工 | `internal/panel/git.go` 整支＋`composer.go:246`＋`cmd/wisp/panel_pump.go` 装配 |
| A3-未提交 | 未提交枚数 | **票 189**（与 diff 同堆，同一枚 struct 加字段） | 189 checked=0／unchecked=6；设计核件 `docs/evidence/s1/189-uncommitted-count-design-a1.md` 已在册 |
| A3-未推 | 未推枚数／领先落后 | **⚠ 无人认领，须立新票或摆 owner**（`docs` 级零要求、Go 侧零源、`GitView` 零这一维） | 尺见 §2.1 A3；先例：`git.go:113` 逐字"remoteBranches is deliberately NOT a field" |
| A4 | 上下文那一行 | **票 214（附件接线）＋票 167（占用）＋票 92（键集）** | 214 Status＝**待派**（0／0 枚 AC 框＝票面还没写 AC 框）；167 checked=0／unchecked=7 |
| A5 | 活动／时间线 | **票 145**（十四行表行 1–5 与行 9–11 那一族；不是右栏另立一堆） | `PLAN.md:3481-3494`＝那张表的**14 枚数据行**，行 12＝成本 |
| A6 | 审批队列 | **票 145 行 6/7 已射程内 ＋ 票 219（三按钮＋理由框）** | 票面：`PLAN.md:3486`／`:3487`；219 在册 |
| A7 | 文件树 | **票 190（`AC#1` 已勾＝设计核已交）** | 190 checked=1／unchecked=5；设计件 `docs/evidence/s1/190-file-tree-design-a1.md` |
| A8 | 插件／技能清单 | **票 215（`215a` 插件／`215b` 技能）** ＋ 票 168（浏览器插件是另一枚东西） | 215 Status＝**待派**；⚠ `215b` 要动 D36 配置树＝契约面 |
| A9 | 代码审查 diff | **票 189** | 同 A3-未提交 |
| A10 | 内置终端 | **地基＝票 163（立而不派）；面板载体＝票 191（`AC#0` 已裁"做"）** | 163 checked=1／unchecked=5；191 checked=1／unchecked=6 |
| A11 | 浏览器／产物预览 | **票 192** | `AC#0` 已按"零 key 抓结果页"定案（票面 Status 行逐字） |
| A12 | 成本与 token | **票 145 行 12（`AC#2`／`AC#6` 未勾）＋票 44（C23 计量本体）** | `PLAN.md:3492`；145 checked=4／unchecked=2＋`AC#2b` |

⇒ **归并结果：12 堆全部有归属指认**（落在 145／163／164／167／181／188／189／190／191／192／214／215／219 这 **13 枚**既有票上；A5「工具调用行」我也已指回票 145，**不新立票**）。**但拆到子项，仍有一枚无人认领＝`A3-未推枚数`** ⇒ **本程真正要点名立票／摆 owner 的＝1 处**。

⚠ 与上一轮的分歧：`182-c1` 写"没票认领＝7 处要点名立票"、`182-a1` 写"候选池无人认领＝K5 插件与 K7 的面板载体与 K10 浏览器"——**到今天这两句都不成立了**：K5→票 215、K7 载体→票 191、K10→票 192 全都已经立票且部分已裁（见票面 Status 行）。§7 P5 具名。

### 3.2 特别那一问：面板要的"宿主侧读面"与模型侧工具（`fs.list`／`shell.exec`）要不要共享实现？

**票面要求写清并交编排者裁。本程给两案的代价、不自选：**

- **共享**（面板直接复用 `internal/tools` 那七枚 `fs.*`）：省一枚读面，但把 **D34 的风险分级渗进面板通道**——`fs.list` 的越界语义是"判 L2 交人批"（`internal/risk/rules_gateway.go` 那一族），而面板今天**没有任何一枚批准入口**（`internal/panel/composer_dispatch.go:176-202` 只答 6 枚方法，其中没有 `decide`／`allow`）。共享之后只有两条路：要么面板读越界时**静默少读**（＝票 92／145 那条"空数组不许顶替"的同族坑），要么面板被迫长出一枚批准面（＝撞 `AGENTS §1.2` 那句"由面板侧来源的 L2『允许』"）。
- **不共享**（各自一枚读面）：代价是**两包各拼一份路径解析**（`A415` 那一族老病），且第二份必须同样过 C26——而这条已经有一份**先例可抄**：票 181 的宿主侧 `internal/panel/git.go:164-172` 自陈"nothing here resolves a path through C26 and nothing here decides permissions"，它读的是**宿主已经被指过去的那棵树**。
- **本程能证的、只到这一步**：票 190 的 `AC#1` 定案已选"**树腿不接受面板传根**、根＝配置根∩工作区"（票面逐字），这一句在语义上已经把 A7 推向"不共享"那一侧；而 A10 终端那一堆**票 191 已把方向定了**：终端不执行命令，它把命令交给既有 `shell.*` 那一腿走**同一套门控与审计**（票 191 面逐字），＝那一堆**是共享**。**⇒ 两堆今天的定案方向相反**，这不是我能替裁的，交编排者。

---

## 4. AC#4 原文：与票 181 的分工写死

| 线 | 归谁 | 今天真身 |
|---|---|---|
| git **状态显示**（当前分支／detached／工作树清单／本地分支清单／`kind`） | **票 181** | **已落地**：`internal/panel/git.go:117-145`（10 枚字段：`kind`／`reason`／`branch`／`detachedSha`／`isDetached`／`repoRoot`／`currentWorktree`／`worktrees[]`／`branches[]`／`switchBlocked`）＋`composer.go:246`＋泵装配 |
| git **内容审查**（逐文件 diff、改了哪几枚文件、未提交枚数） | **本票那一堆＝票 189 落地** | **零字节存在**：尺 0 命中（§2.1 A3），`GitView` 里零这一维 |

⚠ 分工落法（避免"归口充数"）：**票 181 的 AC 框本程一枚没碰**；两票的接口是**同一枚 struct 加字段**（`GitView`）**而不是加快照键**——这条有在册先例，`git.go:113` 逐字"remoteBranches is deliberately NOT a field"，同 struct 加字段不动顶层键集，才避开 `Q-51` 那枚双向对账。**这一句是设计约束，不是我的裁定**。

### 4.1 "未提交枚数"这一枚在 Go 侧到底有没有源

**没有。三个层次分别答清：**

1. **真值产出者：无。** 尺（`cmd internal tools` 三根全给，排测试）：`grep -rniE "git status|porcelain|numstat|git diff"` ⇒ **0 命中**；正控同形尺 `grep -rniE "ReadGit"` ⇒ **10 命中** ⇒ 尺能响、0 是真 0。
2. **现成的 git 取数点：全仓只有一处，且不是 status。** `tools/d22scan/gitignore.go:373 cmd := exec.CommandContext(ctx, "git", args...)`，两发子命令逐字是 `ls-files -z`（`:274`）与 `ls-files -z -i -c --exclude-standard`（`:278`）＝**跟踪清单／忽略匹配**，不是工作树脏净。且 `tools/d22scan/go.mod` 是**独立 module**（`module github.com/CarlosShao/wisp/tools/d22scan`），对 `internal/**`＋`cmd/**` **零 import** ⇒ **不是一枚可复用的读面，是一台 CI 仪器**。
3. **载体：无。** `GitView` 十枚字段逐枚名过（§4 表），无 `uncommitted`／`files[]`／`diff`。

**本程现量（带时刻，`2026-10-03T10:37:02+08:00`，同一台机器、同一棵共享树）：**

| 口径 | 尺 | 读数 | 说明 |
|---|---|---|---|
| **整棵树** | `git status --porcelain \| wc -l` | **473** | 分档：` M`16／` D`16／`??`441（`cut -c1-2 \| sort \| uniq -c`） |
| **只算产码根** | `git status --porcelain -- internal/ cmd/ \| wc -l` | **2** | 两枚都属**在飞的 `255-r2` 写腿**，本程名册差集为空 |
| **未推枚数** | `git rev-list --left-right --count cnb/dev...HEAD` | **115**（behind 0） | ⚠ **会随并发漂**：三枚写腿在飞，任何引用必须带时刻；且这**是"整棵树相对 `cnb/dev`"，不是"我这轮"** |

⚠ **"未提交"与"未推"是两枚不同的洞**：票面 `AC#1` 把"未提交枚数"写在**环境信息**那一堆、"提交或推送"写在同一行；两枚在 Go 侧**都是零源**，但**只有"未提交枚数"有人认领（票 189）**，"未推枚数"**无人认领**（§3.1）。另：`189-uncommitted-count-design-a1.md` §3 已把"暂存 vs 未暂存算不算一起"摆成 owner 的口径题（逐字"＝owner 的口径问题，不是我的"）⇒ 本程不复算、不替裁。

### 4.2 ⚠ 同一枚数在**本程内 9 分钟**就漂了（这是"任何引用必须带时刻"的现证，不是推论）

| 尺 | `10:37:02+08:00` | `10:46:39+08:00`（交件前复量） | 漂了多少／为什么 |
|---|---|---|---|
| `git status --porcelain -- internal/ cmd/` | **2 枚**（`M cmd/wisp/config_reload.go`＋`?? cmd/wisp/restart_tier_keys_255r2_test.go`） | **1 枚**（`?? internal/tools/paths_twocontainments_252_r2_test.go`） | 前两枚被编排者代提交掉（`67ab595d 255-r2（编排者代提…）`），换成一枚**另一条写腿**（`252-r2`）的件 ⇒ **"整棵树的未提交名册"今天在这九分钟里换过一次主** |
| `git rev-list --left-right --count cnb/dev...HEAD` | `0  115` | `0  120` | **＋5 枚**＝`git log --oneline e16e3037..HEAD` 现量那五枚，逐枚名：`67ab595d`（255-r2 编排者代提）／`1289d8bb`（probes/252-r2 骨架）／`3bc10b65`（220-a2 骨架）／`17cfc416`（ledger A565 收 220-r1）／`bfdd910a`（145-a2 骨架）⇒ **一枚都不是本程的**（本程此刻尚未 commit） |
| 独立旁证（同一枚尺、别人手里） | 票 228 收表那枚 commit `e16e3037` 的正文逐字记「未推 112 枚」 | 本程两读数＝115／120 | ⇒ **112 → 115 → 120 三个时刻串得起来**，"未推枚数"这一维**每收一次腿就涨**；任何把它画进界面的票，字段语义必须先把"谁的时刻／哪棵树相对哪个远端"钉死 |

⇒ **本程结论的引用规矩（写给下一位）**：凡引本 §4.1 那三行读数，**连时刻一起引**；凡"未提交枚数"落地，**必须同时答"分母是整棵树还是本次会话那一族文件"**——这一格今天在本仓**没有任何一处文档写过**（票 189 面只写了"口径＝owner 的"）。

---

## 5. AC#5 原文：雷区清单（逐堆标"有没有由面板发起的批准／写动作"）

### 5.1 今天的 inbound 名册真身（这是判雷区的底尺，现读）

- **答得上的方法共 6 枚**：`internal/panel/bridge.go:42-45` 四枚 `panel.*`（`mode.request`／`workspace.request`／`attachment.add`／`message.send`）＋ `:66-67` 两枚 **无前缀**（`config.get`／`config.set`，票 248）。路由逐枚在 `composer_dispatch.go:177-202` 六个 `case`。
- ⚠ **那枚"恰四枚 `panel.*`"锚钉看不见后两枚**：`internal/panel/git_test.go:407 whitelistMethodsFromSource` 只按 `panelMethodRe` 抽 `"panel.` 字面量，`:385`／`:517` 两处钉的是**恰好 4**（`:514` 的正控是"五枚版本必须读出 5"）。⇒ **今天那两枚无前缀写方法不占这枚锚的数**，`bridge.go:54-56` 逐字承认这是**刻意的**（"They carry no `panel.` prefix on purpose, and that is load-bearing twice over"）。**结论：锚钉能挡"第五枚 `panel.*`"，挡不住"第六、第七枚无前缀写面"。** 这一格是要人知道的现状，不是本程要修的缺陷。
- **批准面：零。** 能力尺（**按函数名判、不按词面判**）：`grep -rniE "func .*(allow|decide|veto)" --include=*.go internal/panel` 非测试＝**0 枚**；同尺的**正控打在 `internal/tools`** ⇒ **3 枚**（`cancel.go:63 Vetoed`、`paths.go:53 NewPathCanonicalizer`、`paths.go:189 InAllowlist`）⇒ 尺能响、0 是真 0。⚠ 我**先试过词面尺** `grep -rniE "\bdecide\b|\ballow\b|approve"`（非测试＝**67 命中**），**那把尺不作判据**并已弃用：67 命中里绝大多数是注释里的禁令句子与 `auto_approve` 这枚模式名（`composer.go:155 L2ConfirmTarget = risk.ModeAutoApproveName`），它数的是"提没提到这个词"、不是"有没有这条路"。规格侧同向：`PLAN.md:3487`（§17.5 行 7"面板内 L2 卡"）逐字"按钮区只有 **`拒绝`（次级）** 与 **`查看完整参数`（ghost）** —— **没有"允许"按钮**（D33/F2 定案）"。

### 5.2 逐堆雷区判定（⚠＝这一堆里存在／会被要求存在面板侧动作）

| # | 堆 | 有面板发起的批准／写动作？ | 具名指回 |
|---|---|---|---|
| A1 | 子代理 | ⚠ **半枚：停止** | `task.cancel` 今天是**模型侧工具**（`internal/tools/task.go:625`），面板要"点一下停掉某枚子代理"＝**新增 inbound 方法**＝C17 契约面；`git_test.go:517` 那枚锚会红（若带 `panel.` 前缀）。197 面逐字"nothing here claims a subagent can approve its own (ticket 197 AC#5)" |
| A2 | 后台任务 | 同上，同一枚 | 同上 |
| A3 | 环境信息 | ⚠ **两枚写动作：提交／推送** | 票 181 面 §本票不解决 第 1 条：`Q-64` 只裁了"切分支／切工作树"（→票 186），**提交／推送从来不在批准里**；两者都要**新 C17 方法** ⇒ 契约级、要编排者先落批准记录 |
| A4 | 上下文那一行 | ⚠ **一枚：附件移除** | `panel.attachment.add` 在册（`:44`），**`attachment.remove` 不存在**（名册尺：6 枚方法里没有它）。删除＝写盘动作 ⇒ 与票 214 同堆，别当"顺手加个按钮" |
| A5 | 活动 | 否（纯显示） | — |
| A6 | 审批队列 | ⚠⚠ **这一堆是整个雷区的中心** | `PLAN.md:3487` 逐字"❌ 面板上放'允许'按钮（已被 §15 第 6 项定案否决）"；**票 219 标题逐字"approval card three reply buttons and a reason box"＝正面撞这一格** ⇒ 那是 `Q-49` 那一族＋`AGENTS §1.2`，**契约级、待人拍板**，不在本票射程 |
| A7 | 文件树 | ⚠ **两枚：删除／移动** | 模型侧已有 `fs.delete`（`fs_write.go:583`）／`fs.move`（`:440`）＝**受门控的 Tool**；把它们做成面板按钮＝把 `AGENTS §1.2` 那句"把宿主内部 artifacts 写入实现成受门控的 Tool"与"面板侧来源的 L2『允许』"两条禁令一起撞 |
| A8 | 插件／技能 | ⚠⚠ **最硬的一枚：插件开关** | 上一轮逐字（`182-a1`，转记）："条目带 `on:`／`level:`＝写配置＋改档位"。本程现量补一刀：**今天这条路是关着的**——`internal/panel/config_handlers.go:108 lockedFieldFamilies` 把 `"plugins."` **列进锁死族**，`refusedLockedFamily`（`:115-126`）默认理由逐字"该段属于安全锁定族，放宽必须走带原生二次确认的那条路，不走设置页" ⇒ **`config.set` 不认插件段**（`risk.`／`fs.`／`net.`／`privacy.`／`models.`／`audio.` 同锁）。开关真要能拨，落点是**新面**、不是复用 `config.set` |
| A9 | 审查 diff | ⚠ **两枚：接受／回滚** | 转记 `182-c1`／`182-a1` §12；本程复算：Go 侧连 diff 的读面都没有（§4.1），"接受／回滚"＝`git apply`／`checkout` 一族＝**面板侧写宿主树**，`AGENTS §1.2` 直邻 |
| A10 | 内置终端 | ⚠⚠ **执行面进面板通道＝这次最大的一枚** | 票 191 面已把它写成"为什么它不能像其他几堆那样直接补"，逐字指回 `PLAN.md:2027` ①③ 与 R20；`AC#0` 已裁"甲"（做），**裁的是做法、不是放宽边界** ⇒ 落地必须走"面板不执行、把命令交给 `shell.*` 走同一套门控与审计"那一支（票 191 面逐字） |
| A11 | 浏览器／产物预览 | ⚠ **一枚：外泄／注入管子** | 预览本地 HTML＝在面板里渲染宿主生成的内容＝票 143／175／177／183 那一族的同一条管子；⚠ 本程量到 `web.*` 零实现，所以这条雷**现在踩不到、一旦补源就立刻在** |
| A12 | 成本 | 否（纯显示） | — |

⇒ **⚠ 堆数现量：12 堆里有 9 堆带雷**（A1 A3 A4 A6 A7 A8 A9 A10 A11），其中**契约级、必须编排者先落批准记录**的＝**5 枚**：面板批准（A6／票 219）、插件开关（A8）、终端执行面（A10）、git 提交／推送（A3）、diff 接受／回滚（A9）。**本程没提"顺手把批准接上"，九枚一律单列。**

---

## 6. 量不到的格子（本程如实登记，不猜）

| # | 量不到的那一格 | 为什么量不到 | 要谁才量得到 |
|---|---|---|---|
| M1 | A12 成本：往 `PumpSources` 补一枚 `Cost` 读口、往快照加一枚键，会不会红那三枚键集／双向对账尺 | ⛔ 任务书 §2 禁跑 Go 命令（三枚写腿在飞）。我只能证**今天没有读口**，证不了**加了会不会红** | 一枚能跑 `go test ./internal/panel/` 的程 |
| M2 | A5 活动：`turn.ToolCalls` 有没有一条**不改 loop 形状**就能被泵读到的路 | 同上；且 `145-c2` 已把这条写成"数据只在被打印的那一刻存在"（`docs/reports/` 那句，我没复算） | 票 145 的落地腿 |
| M3 | **这一栏界面上本来打算显示哪几枚子行**（A4 三行、A5 时间线、A8 分组标题） | ⛔ `design/**` 与 `frontend/**` 两层禁读 ⇒ 上一轮的主尺（`RbPanels.register`）本程不可复跑，我只得**转记**、不得自证 | 能读 `design/doubao/demo/**` 的那一枚会话（编排者口径：那一侧归别家） |
| M4 | 12 枚堆与**产品真身**的对应（哪些堆其实已经被 owner 划到 composer 那一行、不属于右栏） | 同 M3——堆的"归属哪一屏"只有 demo／前端树能证 | 同 M3 |
| M5 | A1／A2 的行内字段在**页面上画成什么**、`tasks` 键前端声明了没有 | ⛔ `frontend/src/lib/panel.ts` 禁读。只量到 Go 侧自陈：`composer.go:87-90` 逐字承认那枚双向尺现在**点名 `tasks` 是 interface 没声明的键** | 别家会话（`Q-51`） |
| M6 | 未提交枚数的"暂存 vs 未暂存"口径 | 不是量的问题——`189-a1` 已写成 owner 的口径题；本程不替裁 | owner |
| M7 | "未推枚数"该画成谁的数（整棵树 vs 本次会话） | 不是跑不跑得的问题——**本仓没有任何一处文档定义过这一维的口径**（票 189 面只写了"暂存 vs 未暂存＝owner 的口径题"，未推那一维连这句都没有）。本程**已用两把同尺现证它会漂**：9 分钟内 `115 → 120`、脏名册换主（见 §4.2）⇒ 落地前必须有人裁口径，不是量一次就完 | 编排者或 owner（要落进票 189 或新立一枚） |

---

## 7. 推翻清单（票面与派单前提，逐条具名）

| # | 被推翻的那句 | 逐字出处 | 本程真身与尺 |
|---|---|---|---|
| **P1** | "Go 侧'子代理'实体 **0 命中**" | 票 197 现量表第一行（尺逐字 `grep -rln "subagent\|Subagent\|SubAgent" --include=*.go internal/ cmd/ tools/ \| grep -v _test.go \| wc -l`） | **非测试命中文件＝9 枚**（`internal/panel/composer.go`／`pump.go`／`subagent_roster_197.go`、`internal/streamkey/streamkey.go`、`internal/tools/bridge.go`／`subagent_197.go`／`task.go`、`cmd/wisp/panel_pump.go`／`run.go`）。**性质＝过期，不是错**：那一句量在票 197 **落地之前**，票 197 自己就是补它的那一程。⚠ 票面原文**不改**，本程只追加 |
| **P2** | "`TaskRoster` ……**没有'状态'这一维**" | 票 182 现量第 4 条（`182-c1` 段）与票 188 现量表第二行 | **状态维已经在**：`internal/tools/task.go:158 func (o TaskOutput) StateAnswer() (statemachine.State, string)`，`:159-165` 两形 fail-closed（没登记→"宿主没有登记这一维"；登了但不是 D43 那 20 枚名字之一→点名"票 188 AC#2：不许自造"）。⇒ 同一句"名册在、状态维没有"**今天不成立**；票 188 `AC#2` 的射程随之前移 |
| **P3** | "内置终端 源＝`shell.session`……该工具**今天未注册**" | 票 182 现量第 4 条第二子条 | **比"未注册"更黑：实现体也不存在**。尺 `"shell.(exec|session)"` 非测试 **0 命中**（正控 `"task.(spawn|output|cancel)"`＝6）；`internal/tools/` 非测试件名册里**没有 `shell*.go`**（`bridge/cancel/capability/doc/fs/fs_edit/fs_staging/fs_write/gate/grant/mode/paths/paths_workspace/platform_other/recycle_windows/registry/subagent_197/task/task_backfill/tool/volume_windows`）。⇒ 票 163 的地基账不是"接上线"，是"从 `shell.exec` 那枚 1MB 上限（`PLAN.md:432`）起做整枚工具"。**这条直接改 A10 的代价档** |
| **P4** | 派单 §0／票 182 标题那句"这一栏**整块**都不在票 145 那 14 枚状态的射程里" | 派单 §0 原话＋票 182 标题 | **不成立（过强）**：14 枚射程＝`PLAN.md:3481-3494`（§17.5 那张表的 14 枚数据行），本栏的 **A6 审批**（行 6 `:3486`／行 7 `:3487`）与 **A12 成本**（行 12 `:3492`）**就在射程里**——`182-c1` 自己就是这么归的（逐字"可并格只有 S7→票 145 行 12"）。**准确说法＝12 枚堆里 10 枚在 145 那张表之外、2 枚在其内**。**这一条对本票的实际后果**：派单 §0 那句"要么没源、要么有源没载体、要么规格里根本没人要求过"仍对，但**"没人要求过"那一档不能整栏套用** |
| **P5** | "归属：没票认领＝**7 处要点名立票**"（`182-c1`）／"候选池无人认领、须立新票＝**K5 插件与 K7 的面板载体与 K10 浏览器**"（`182-a1`） | 票 182 Progress log 两段 | **都已不成立**：K5→**票 215 已立**（Status 待派，且已拆 `215a`／`215b`）、K7 载体→**票 191 已立且 `AC#0` 已裁**、K10→**票 192 已立且 `AC#0` 已定案**、文件树→**票 190 `AC#1` 已勾**、审查→**票 189 已立**。⇒ 本程复算后**真正常无人认领的只剩 1 处：未推枚数**（§3.1）。⚠ 前两程当时没算错，是**它们之后新立了 6 枚票**（189–192／214–215）；引用这两句时必须带票号 |
| **P6** | 编排者 09-28 17:5x 钉死的口径："本栏堆数＝**原型 `design/doubao/demo/**` 注册表口径**" | 票 182 Progress log 末段 | **本程无法执行这条口径**（任务书 §2 把 `design/**` 一并禁读）。⇒ 本程交的是**另一把尺的数（12 枚，可点名数据堆口径）**，两把尺并存登记、不互相顶替。**要么编排者解禁 `design/**` 给下一程，要么明写本栏堆数改口径**——本程不自选（这就是"未定义即停"那一格） |
| **P7** | `182-a1` 段末那句"现读 `grep ... "type Record struct\|type Job struct\|type BackgroundJob"`＝0 命中，故派单 §E 说的状态维写面在 `internal/memory/**` 不成立" | 票 182 Progress log `182-a1` 段 | **那一记推翻到今天是真**（`TaskRoster`／`TaskOutput` 仍是真身，`internal/tools/task.go:176`／`:158`），但它的**结论已过期**：它当时推的是"状态维不存在"，现在状态维存在（P2）。⇒ 引用那一格要连着 P2 一起引，不然下一位会拿它当"仍然没有状态维"的证据 |
| **P8** | 票 182 现量第 3 条："`Snapshot` 四枚（`internal/panel/composer.go:57-62`）＋`ComposerState` 六枚（`:200-207`）" | 票面现量第 3 条 | **两枚数都涨了**：`Snapshot` **6 枚**（`composer.go:57-92`：`pending`／`results`／`composer`／`generatedAt`／`instructions`／`tasks`）、`ComposerState` **11 枚**（`:235-270`：加 `git`／`currentModel`／`modelKnown`／`credentialState`／`credentialKnown`）。⚠ 这一格 `145-c2`（今天 09:5x，锚 `5f9ff9d4`）已经先量到并写进票 145 ⇒ **本程是复算同值、不是新推翻**，写在这里是为了让票 182 面的那两行不再被照抄 |
| **P9** | `PLAN.md` 的行号引用（票 182 现量第 2 条尺 `sed -n '3473,3488p' docs/PLAN.md`；票 145 头行"`PLAN.md:3473`（表头）／`:3475-3488`（十四行数据）"） | 两枚票面 | **行号已漂 6 行**：今天表头＝`PLAN.md:3479`、分隔行 `:3480`、**14 枚数据行＝`:3481-3494`**。⇒ 两枚票面那把尺今天**复跑读到的是错的那 6 行**（会漏掉行 13／行 14 两枚、并吞进 §17.5 的标题与导语）。⚠ 内容本身没漂：D34 行 `:2568`／`:2569`、REJECTED 多 Agent `:1544`、并发闸门 `:1447`、§17.5 行 12 成本 `:3492` 本程逐枚现读**同值**。⇒ 复跑任何 `PLAN.md:` 行号尺之前先按节名定位，别照抄行号 |
| **P10** | 票 176 标题逐字"There is no spawner: `RunAsync` has **zero production call sites** and `TaskRoster.Record` has **zero writers**"（这格是 A2 那一堆的上游前提） | `.scratch/wisp/issues/176-...md` 标题 | **两句今天都不成立**：`RunAsync` 生产调用者＝**2 枚**（`cmd/wisp/run.go:1099`、`internal/tools/subagent_197.go:344`）；`TaskRoster.Record` 写者＝**至少 2 路**（`internal/tools/task.go:343` 经 `PublishSubagent`、`internal/tools/subagent_197.go:424` 直接调）。⇒ ⚠ **这一格不在本票射程内、我没翻它的意思**：票 176 面 AC 框一枚没碰，本程只登记"它的现量前提已被票 188／197 取代＝过期非错"，供编排者决定是否在 176 面追加一节。**顺带一条要写腿注意的现量**：`internal/tools/task_backfill.go:12` 的注释**今天仍逐字写着** "(Loop.RunAsync, loop.go:321) with zero production call sites"＝一枚**已与代码对不上号的注释**，它会误导下一位对 A2 的判断（尺：`grep -rn "zero production call sites" --include=*.go internal cmd`） |

---

## 8. 三档定性表（交件硬要求·汇总视图）

| 档 | 枚数（整堆＋半堆，与 §2.2 同算式） | 哪几堆 | 判据一句话 |
|---|---|---|---|
| **有源**（缺载体或缺投递） | **5 整堆＋1 半堆** | A1 子代理（三处齐）、A2 后台任务（三处齐）、A6 审批（三处齐）、A5 活动（文本流三处齐／**工具调用行**缺）、A12 成本（缺读口＋缺键）＋ **A3 的树那一半**（三处齐） | 三问（产出／投递／落盘）逐问有名 |
| **无源但规格要求** | **2 整堆＋1 半堆** | A10 终端（**D34 `PLAN.md:2569` 行在册**，实现体零）、A9 审查 diff ＋ **A11 的网页那一半**（**D34 `:2546-2548` 三行在册**＋`AGENTS §2` 未定案） | 要动的节：A10／A11＝**动实现不动契约**（D34 行已有）；A9＝**只到票面级**（`docs/specs/**` 对 diff 零要求，从严即落"真空"，本程从宽并注） |
| **规格真空** | **2 整堆＋1 半堆** | A7 文件树、A8 插件／技能清单 ＋ **A11 的产物预览那一半** | `docs/specs/**`＋`PLAN.md` 对该维**零要求**（R-C 尺）；做＝新需求（多数已被票面接单） |
| **第四种形状：接了半根线** | **1 整堆** | A4 上下文那一行 | 快照键集**在**、`pump.go:269` 那一跳把附件参数**硬写 `nil`**、broker 生产零调用者、预算子维无源 ⇒ 上面三档都装不下它 |
| **A3 余下两子项（半堆之下再拆一层）** | 2 枚子项 | 未提交枚数（无源、**票 189 认领**）、未推枚数（无源、**无人认领**） | 不许整堆归一档；见 §4.1／§4.2 |

**加总校验（把半堆还原成整堆）**：三档里被拆开的只有 **A3**（树那一半＋未提交＋未推＝1 整堆）与 **A11**（网页半＋预览半＝1 整堆）⇒ 12 枚＝**5＋2＋2＋1＝10 枚整堆不拆** ＋ **A3** ＋ **A11**，无重复计数、无遗漏。

**契约级判语（哪几档的落地会撞 `C1–C32`／`D1–D47`）**：A10／A11 **不撞**（D34 行已在册，补实现即可，⚠ 但 `shell.session` 一旦注册＝D34 那行从"文档"变"运行时"，票 172 那枚"把冻结档位钉到门实际给的档位"的尺会开始响）；A8 的 `215b` 技能那一半**撞 D36**（配置树，票 215 面逐字）；§5.2 那 5 枚契约级雷**全撞 C17**（白名单加方法）。

---

## 9. 台件与尺台账（复跑用）

本件全部命中数出自以下尺，逐把带正控；根一律显式写死（`cmd internal tools docs scripts` ＋票池），无一把用 `.` 当根：

1. `grep -rniE "git status|porcelain|numstat|git diff" --include=*.go cmd internal tools`（排 `_test.go`）＝**0**；正控 `grep -rniE "ReadGit" --include=*.go cmd internal`（排测试）＝**10**。
2. `grep -rn "os.ReadDir|filepath.WalkDir" --include=*.go internal/panel cmd/wisp`（排测试）＝**3**（逐枚列出，见 §2.1 A7）。
3. `grep -rnE "\"shell\.(exec|session)\"" --include=*.go cmd internal`（排测试）＝**0**；正控 `grep -rnE "\"task\.(spawn|output|cancel)\"" --include=*.go cmd internal`（排测试）＝**6**。
4. `grep -rln "ubagent" --include=*.go internal cmd`（排测试）＝**9 枚文件**（逐名列见 §2.1 A1）。
5. `grep -rniE "plugin" --include=*.go internal/panel cmd/wisp`（排测试）＝**3**（逐枚：`config_handlers.go:108` 名册一行＋两枚注释）；正控 `grep -c -i plugin internal/config/schema.go`＝**18**、`grep -c -i skill internal/config/schema.go`＝**0**。
6. `grep -rniE "\bCost\b|price|micros" --include=*.go internal/panel cmd/wisp/panel_pump.go`（排测试）＝**13**，逐枚看过无一枚是运行成本（正控＝同尺 `\bModel\b`＝**40**）。
7. `grep -rniE "occupancy|ContextBudget|BudgetUsed" --include=*.go internal/agent internal/panel cmd/wisp`（排测试）＝**0**（正控同上 `\bModel\b`＝40）。
8. `grep -rn "NewAttachmentBroker" --include=*.go cmd internal tools docs scripts` ⇒ 非测试**只有定义那一行**（`internal/panel/attachments.go:179`）；`wc -l internal/panel/attachments.go`＝**410**。
9. `grep -rn "d22scan" --include=*.go internal cmd`（排测试）＝命中全为注释；`tools/d22scan/gitignore.go:373` 是全仓**唯一**非测试 git 起进程点。
10. `grep -rniE "文件树|内置终端|任务监控|子代理|产物预览" docs/specs/ docs/PLAN.md`＝`docs/specs/**` **1**／`docs/PLAN.md` **4**（逐枚见 §1.1 R-C）。
11. `git status --porcelain -- internal/ cmd/`（`10:37:02+08:00`）＝**2**；`git status --porcelain`（全树）＝**473**；`git rev-list --left-right --count cnb/dev...HEAD`＝**0 / 115**。
12. `grep -n "panel\." internal/panel/bridge.go`＝常量 4 枚 `panel.*`＋2 枚无前缀；路由 `grep -n "case Method" internal/panel/composer_dispatch.go`＝**6 枚 case**。
13. 批准面能力尺 `grep -rniE "func .*(allow|decide|veto)" --include=*.go internal/panel`（排测试）＝**0**；正控同尺换根到 `internal/tools`＝**3**（`cancel.go:63`／`paths.go:53`／`paths.go:189`）。⚠ 词面尺 `\bdecide\b|\ballow\b|approve`（非测试＝67）**已弃用**，理由见 §5.1。
14. 工具实现体尺 `grep -rn 'func (.*) Name() string { return "web' --include=*.go internal cmd`＝**0**；字面尺 `"web.`（排测试）＝**3**，逐枚＝`internal/risk/provenance.go:88`／`:120`／`:137`。

---

## 10. 纪律面自陈（本程）

- **AC 框：票 182 与任何别家的勾，一枚没碰、没翻。**
- **产码：零字节。** 写面只有 `.scratch/wisp/probes/182/a1/census.md`。`docs/**`（含 `PLAN.md`／`specs/**`／`evidence/**`）**零字节改动**；`internal/**`／`cmd/**`／`tools/**`／`frontend/**`／`design/**` **零字节改动**。
- **Go 命令：一枚没跑**（`go test`／`go build`／`go vet`／`go list`／`wisp` 全无）——三枚写腿在飞（`255-r2`／`171-r3`／`252-r2`），§6 那 7 格因此只能记"量不到"。
- **`frontend/**`／`design/**`：不读、不引、不转述内容**（唯一例外＝按票面文本**转记**上一程已经量过的读数，逐处标"转记"并给 `docs/evidence/s1/...` 出处）。
- **零删除命令**；未 push；未 `git add -A`／`--amend`／`reset`／`stash`／`checkout .`／`clean`；commit 全部带显式 pathspec。
- **AC#7 排程那一格（票面要求，本程答一句）**：本件是**纯读＋一枚探针件**，不动任何产码 ⇒ **不阻塞**票 179／176／177／175 那条"超长输出读回来"的主干（且 §7 P10 只是登记那枚前提过期，**没替它们翻任何勾**）；也**不构成**"界面先做、后端后补"的理由——本件 §2 的三档就是把"后端有没有源"先数清楚。前端那一侧渲染**未读、未引、未转述内容**（只按票面文本转记上一程已量过的读数，逐处标"转记"）。
- **占位符自查**：本件写完用一把"空缺字样"尺复扫自身（五种常见空缺词形＋中文"待／留"两形），首轮**命中 1 枚＝那把尺自己的字面文本写进了本件**（本行就是它的更正位），改成只描述尺、不写尺面之后 ⇒ **复扫命中 0**。尺面字串刻意不落在此件内（不落＝不落进被扫文件，这是任务书 §5 那条自查的本意）。
- 起手名册与终态名册差集：本程只新增自己这一枚探针件，`internal/ cmd/` 两枚在飞件（`255-r2`）**未碰、未提交、未撤销**。
