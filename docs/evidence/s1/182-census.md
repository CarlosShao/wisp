# 182 — 右上角「任务监控」栏 10 堆的 Go 侧源 ＋ 逐堆归属票（只读普查，锚 `232b0846c`）

> 票 182，AC#2–AC#7。**只读腿**：`internal/**` `cmd/**` `docs/PLAN.md` `docs/specs/**` 与全部票面**一字节未动**，
> 零产码、零 `go` 命令、零 `git` 写（只有 `git log` / `git status` / `rev-parse` 这类读），`frontend/**` **未读未引**（按 AC#7 从严）。
> **写面只有这一枚文件**：`docs/evidence/s1/182-census.md`（新建）。
> 锚：起手 `232b0846c`（`git rev-parse --short HEAD` 起手现跑）；跑到一半被推进到 **`d1b1a2d22`**（他人 `323-ac5` 一笔，
> 只碰 `internal/config/settings_320_test.go` 与 `docs/evidence/s1/323-ac5.md`，**不在本程任何一把尺的射程内**）。
> 收尾复跑 `git diff --name-only 232b0846c..HEAD`＝**2 枚**，与本表所有读数无关；下表读数一律是**工作树现读**。
> 每把尺的原文都在表里，可原样复跑；**命中数一律是「行数」**，不是出现次数。

---

## 0. 口径与两件先钉死的事

1. **堆数口径沿用编排者 09-28 17:5x 的裁定**：这一栏的堆数＝原型 `design/doubao/demo/**` 注册表口径，
   **不是产品真身**（产品真身在前端，本票 AC#7 不许我读）。
2. **owner 那张截图是外部产品（qoder）的界面，不是本仓规格**（票 182 AC#1 原文要求明写）——
   下面每一堆的「谁那儿看到的」只指**本仓工作树文件**与**票面/PLAN 文本**，没有一条来自截图。
3. **正控纪律**：凡 0 命中的尺，同一次调用里都打了「已知存在的同形尺」正控；没打正控的读数一律记〔不可判〕（见 §9）。

---

## 1. AC#1 复核：堆数＝10（7 枚有位 ＋ 3 枚有名无位）

| 尺 | 原文 | 读数 |
|---|---|---|
| 主尺（注册） | `grep -rn --include=*.js "RbPanels.register(" design/doubao/demo/` | **7 行** |
| 对照尺（标题） | `grep -rn --include=*.js -E "^\s*label:" design/doubao/demo/` | **8 行**，剔 `app.js:300`（命令面板条目，非竖条堆）＝**7**，与主尺相等 |

**7 枚有位**（demo 工作树现读，`label` 逐字）：

| 位 | 堆 | 注册位 | demo 里的子行 |
|---|---|---|---|
| K1 | 上下文 | `rightbar.js:10` | CWD／附件／上下文预算条 |
| K2 | 活动 | `rightbar.js:85` | 活动时间线（含「审批超时自动拒绝」「fs.delete 被拒 L2」这类带结果的条目） |
| K3 | 审批 | `rightbar.js:118` | 待审批队列＋满额告警 |
| K4 | 文件 | `rb-files.js:10` | 工作区文件树 |
| K5 | 插件 | `rb-plugins.js:8` | 技能／伙伴分组＋`on:` 开关＋`level:` |
| K6 | 审阅 | `rb-review.js:8` | 逐文件 diff＋「Apply 3 changes」 |
| K7 | 终端 | `rb-terminal.js:8` | 标签页＋「新建终端」＋命令输入行 |

**3 枚有名无位**（不在注册表，仓内文本有据；出处沿用 `182-a1` §2.3 与 `issues/77:364/:384`，本程**未复跑**那两处）：

| 位 | 堆 | 仓内出处 | demo 里有位吗 |
|---|---|---|---|
| K8 | 子代理状态 | 票 188 AC#1（owner 已裁「乙：进」）、票 197 | 无（demo 无独立堆） |
| K9 | 后台任务名册 | 票 188 AC#2、票 164 | 无（混在 K2 时间线里） |
| K10 | 浏览器 | `issues/77:364`「审查/终端/浏览器」、`:384` | 无（三标签定案里那枚今天无位） |

⇒ **可点名总堆数＝10**。AC#1 的框 10-08 已由编排者翻勾，本程只复核读数、**不翻任何框**。

---

## 2. AC#2 逐堆三档定性 ＋ 尺与命中数（含正控）

三档定义（票面原文）：**有源**＝Go 里已能取到真值（缺不缺载体另说）／**无源但规格要求**＝Go 取不到但 `PLAN.md`／specs 有人要求过／**规格真空**＝今天没人要求过，做＝新需求。

### 2.1 总表

| 堆 | 三档 | 现读尺（原文）→ 命中数 | 今天载体在哪 | 要动哪一节／是否契约 |
|---|---|---|---|---|
| **K1-A 路径（CWD）** | **有源＋载体已在** | `Workspace WorkspaceView \`json:"workspace"\``＝`internal/panel/composer.go:273`；`awk '/^type ComposerState struct/,/^}/' \| grep -c "json:"`＝**12** | 快照 `composer.workspace` | 不动（票 92 带 `-done`，框态 5 勾／1 未勾） |
| **K1-B 附件** | **有源＋载体已在** | `composer.go:274-277` 四枚（`attachments`／`acceptedAttachmentMimes`／`maxAttachmentBytes`／`attachmentError`） | 快照 `composer.*` | 不动（票 92，同上） |
| **K1-C 占用条（used÷total）** | **无源**（成分口径也没定） | `grep -rn "usedTokens\|UsedTokens\|contextUsed\|ContextUsed" internal/panel --include=*.go \| grep -v _test.go`＝**0**；正控＝同形尺 `json:"mode"\|json:"git"`＝**2**（尺能响） | 无；`PumpSources`（`internal/panel/pump.go:140-204`）11 个 reader 里没有占用 reader | 票 167 AC#2（未勾）；`PLAN.md` 没有分子/分母的定案文字＝**须先定口径** |
| **K1-D token／成本** | **有源缺载体** | 源：`cost.go:52 AddUsage`／`guard.go:148,:155`／`loop.go:523-524,:535-536,:1115 publishUsage`（非测试共 **12 行**）；面板消费：`grep -rn "Usage\|TokensIn\|CostView" internal/panel --include=*.go \| grep -v _test.go`＝**0**（正控同上 2 命中） | **无**：`Snapshot` 6 键（`pending/results/composer/generatedAt/instructions/tasks`）无 cost；`PumpSources` 无 Cost reader | 票 145 **行 12「成本」**，今天位＝`PLAN.md:3510`（旧号 `:3486`）；145 AC#2／AC#6 **未勾** |
| **K2 活动** | **有源缺载体（细维）** | 源：`internal/agent/sink.go:51 type Event`／`:92 RecordSink.Publish`；载体现状＝`ResultChunk{correlationId,text,done}` 三枚，**无工具名/状态/时刻**（`composer.go` ResultChunk 块）；`grep -rn "json:\"feed\"\|ActivityView" internal/panel --include=*.go \| grep -v _test.go`＝**0**；正控＝`json:"rows"`＝**1** | 无 `feed` 键；`PumpSources` 无活动 reader | 内容维**已有规格**：`PLAN.md:3502` 行 4「工具调用」chip 已定义「运行中／成功／失败／被拒 ＋ L0/L1/L2」⇒ 并 **145 行 4／5 的载体扩张**即可；若 owner 要「带时间戳的独立时间线」＝**十四行表新增一行＝改 `PLAN.md` 文字（须人工批）** |
| **K3 审批** | **有源＋载体已在**（序号/停止两维缺） | `ApprovalCardView` 10 个 json 字段（`internal/panel/approval.go:39-59`）；`Position` 在 `internal/panel/**` 非测试 **0 命中**（正控＝`CorrelationID` 在该文件 **3** 命中）；真身：`approval/pending_read.go:58 Position`、`:127` 赋值 | 快照 `pending`；但 `liveVerdicts`（`cmd/wisp/panel_pump.go:64-90`）**逐字段搬运时不带 Position** | 载体格＝145 行 6／7；**序号＝167 AC#3、停止＝167 AC#4**（皆未勾）；不新增键 |
| **K4 文件树** | **有源＋载体已在**（09-28 起最大状态订正之一） | `TreeView` 7 个 json 字段（`internal/panel/tree_types_190.go:50-58`）；`Files *TreeView`＝`composer.go:285`；生产装配 `Files: rt.fileTreeView`＝`cmd/wisp/run.go:922`；读面写调用 `grep -nE "os\.Remove\|os\.Rename\|os\.WriteFile\|os\.Create\|RemoveAll" internal/panel/tree_reader_190.go`＝**0**（正控 `os.Open\|os.Lstat`＝**2**） | 快照 `composer.files` | 不动（票 190 **已结案**，6/6 勾） |
| **K5 插件** | **库里有读面／生产零装配／栏位规格真空** | `func (s *Store) ListPluginStates`＝`internal/memory/dao_misc.go:243`（单测 `dao_test.go:451` 在跑）；**非测试调用者 0**（`grep -rn "ListPluginStates" internal cmd --include=*.go \| grep -v _test.go \| grep -v dao_misc.go`＝**0**）；`UpsertPluginState` 非测试调用者也 **0**（只有定义）；`grep -rn "plugin/manifestslot" internal cmd --include=*.go \| grep -v _test.go`＝**1 行且是注释**（`internal/plugin/command.go:16`）⇒ **生产 import 0**（正控＝`grep -rn "internal/plugin\"" …` 非测试＝**8 行**）；面板侧 `grep -rni "plugin" internal/panel --include=*.go \| grep -v _test`＝**3 行**，逐行过目＝1 行锁定族名单（`config_handlers.go:109`）＋2 行散文，**零读面** | 无 | **无票认领**（全池 `grep -rlniE "插件列表\|已装插件\|右栏插件"`＝只命中票 182 自己）⇒ **立新票或摆 owner**；写动作另见 §6 |
| **K6 审阅（含未提交枚数／diff）** | **有源＋载体已在** | `GitView` **12** 个 json 字段（`internal/panel/git.go:144-182`，含 `Uncommitted`＝`:178`、`Diff`＝`:181`）；进程面：`runGitReadOnly(` 4 处出现（定义 `git_readonly_exec.go:349` ＋ 调用 `:426`／`:493`／`:556`）；命令形状**只有查询**：写动词字面 `grep -nE "\"commit\"\|\"push\"\|\"checkout\"\|\"reset\"\|\"apply\"\|\"restore\"" git_readonly_exec.go`＝**0**，正控查询字面 `"status"\|\"diff\"\|\"rev-parse"`＝**5** | 快照 `composer.git.uncommitted` / `.diff` | 不动（票 189 **已结案**：AC#2–#6 勾，AC#1 未勾但设计件 `189-uncommitted-count-design-a1.md` 与 `189-ac-verdict.md` 都在盘） |
| **K7 终端** | **无源但规格要求** | 注册循环 `grep -c "tools.Builtin.*Entries(" cmd/wisp/run.go`＝**5**（FS `:614`／Task `:643`／System `:670`／Subagent `:1012`／Plan `:1051`，**无 Shell 一支**）；`shell.session\|shell.exec` 非测试＝**6 行全注释**（`config/unwired.go:63/64/70/76`、`risk/rules_shell.go:14`、`agent/spill.go:26`）；PTY 非测试 1 命中＝`cmd/wisp/testdata/esclistener/main.go:117` 的 `ptY` 假阳 | 无 | 契约行在：D34 表 `PLAN.md:2575`「`shell.session`｜在一条活着的会话里跑下一条命令｜**L2**｜票 163」＝**契约行**；地基＝**票 163 AC#1b**（未勾）；载体＝**票 191**（AC#0 已裁**甲**：复用工具腿＋门控＋审计；AC#1–#6 未勾） |
| **K8 子代理** | **有源＋载体已在**（09-28 后第二次订正） | `TaskRowView` **12** 个 json 字段（`internal/panel/subagent_roster_197.go:111-142`，含 `Status`／`StatusKnown`／`StatusReason` 三件套）；`Tasks *TaskRosterSection`＝`composer.go:91`；装配 `Tasks: rt.taskRosterState`＝`cmd/wisp/run.go:943` | 快照 `tasks` | 不动规格；票 **197**（AC#4/#5 勾，AC#0/#1/#2/#3/#6 未勾）＋票 188 AC#1 已裁「乙：进」 |
| **K9 后台任务** | **有源＋载体已在** | `TaskRosterSection` **7** 个 json 字段（`:155-183`：`rows`／`inFlightSlots`／`poolCap`／三个流截断/丢弃读数）；状态词表＝`TaskOutput.StateAnswer`（`internal/tools/task.go:158`） | 同上 `tasks` | 票 **188 AC#2**（未勾，⚠ 前提写「给 `Record` 补状态维」，而 `type Record struct` 在本树 **0 命中**、正控 `type TaskRoster struct` **1** 命中 ⇒ 该格口径已过期，见 §9-3）＋票 164（**已结案**，0 未勾）＋票 196（**已结案**，资料/裁决型无判据框） |
| **K10 浏览器** | **无源＋规格真空**（有票但被前置扣住） | `grep -rli "browser" internal/panel --include=*.go \| grep -v _test`＝**0**（正控＝`grep -rli "WebView" cmd/wisp`＝**14 文件**）；`cmd/wisp` 非测试提 browser 的 **3 文件**（`panel_assets.go`／`panel_control_windows.go`／`panel_host_windows.go`）全是 WebView2 宿主措辞，**没有一枚栏内浏览读面** | 无 | 票 **168**（7 框全未勾；AC#1 前置＝「50/51 未结案不许开工」——现读 **50 已结案、51 未结案**（`ls .scratch/wisp/issues/ \| grep -E "^5[01]-"`：50 带 `-done`、51 不带）⇒ 闸门仍在） |

### 2.2 小结（按"最便宜到最贵"排）

1. **只差载体（全归票 145 那一集）**：K1-D token／成本（行 12）、K2 活动（行 4/5 扩载体）。两枚都是"源已在、字段没进包"。
2. **只差一个字段维（票 167）**：K1-C 占用条、K3 序号、K3 停止。三枚源/半源都在，缺的是定义与载体。
3. **已闭、只须记账**：K1-A／K1-B（92）、K3 主体（145 行 6/7 的 pending）、K4（190）、K6（189＋181 分工）、K8／K9（197／188／164／196）。
4. **无源、票已在手但不许我重开**：K7（163 地基未注册＋191 载体未落地）。
5. **无人认领**：K5 插件（库里有 `ListPluginStates`，生产零装配、栏位零规格、池内零票）；K10 浏览器（168 被 51 扣住）。

---

## 3. AC#3 逐堆归属票（本票主要产出）

> 指认前逐张票现读了 AC 标题与框态（读数见「框态」列）。**本程没有替任何票翻勾**。

| 堆 | 归属（具名） | 并格还是立票 | 该票今天框态（现读） | 一格一句为什么是它 |
|---|---|---|---|---|
| K1-A 路径／K1-B 附件 | **票 92** | 已闭，不并格 | `92-...-done.md`；5 勾／1 未勾（AC#5 变异三向） | `workspace` 与附件四枚都是票 92 落的；未勾那格是变异判据，不影响载体 |
| K1-C 占用条 | **票 167 AC#2** | 并格（该票就是"占用/排队/停止/草稿＋崩溃自救"四枚小输出） | 6 未勾／1 勾 | 占用条的"分子或分母未知 ⇒ 必须显式未知"正是 167 AC#2 的原文 |
| K1-D token／成本 | **票 145 行 12** | 并格（145 的"扩快照"） | 3 未勾（AC#2／AC#2b／AC#6）／4 勾 | 145 的十四行表第 12 行逐字就是"成本（C23）"，五枚字段的源链在 `cost.go`→`guard.go`→`loop.go` 上现读可得 |
| K2 活动 | **票 145 行 4／5**（载体）＋**票 196**（状态词表，已结案） | 并格；**不另立票** | 145 同上；196 已结案（无判据框） | 行 4「工具调用」chip 已把"运行中／成功／失败／被拒＋风险级"定义完了，活动堆缺的就是把 sink 事件搬进 `results`／新细分；状态词表 196 已定案 |
| K2（可选）独立时间线 | **须 owner 一句话** | 若做＝改 `PLAN.md`（十四行表新增一行） | — | 十四行表今天没有"活动流"这一行；加行＝改契约文字，按铁律只摆 owner 批 |
| K3 审批主体 | **票 145 行 6／7**（pending 载体） | 不并（已在） | — | `pending` 已是快照第一维；145 表别重复列一行（防双记账） |
| K3 序号／停止 | **票 167 AC#3／AC#4** | 并格 | 6 未勾／1 勾 | `LiveApproval.Position` 已有真源（`pending_read.go:58/:127`），停止维是 167 AC#4 原文 |
| K4 文件树 | **票 190** | 已闭 | **已结案**，6/6 勾 | `TreeView`＋`files`＋生产装配都在；设计核证据件 `190-file-tree-design-a1.md` 在盘 |
| K5 插件 | **无票认领 ⇒ 立新票或摆 owner（二选一，交编排者）** | 立票 | — | 全池无票覆盖"列出已装插件"；库里 `ListPluginStates` 有读面但零生产调用者，**做栏位＝新需求**（规格真空）；开关那一枚另见 §6 |
| K6 内容审查 | **票 189** | 已闭 | **已结案**；5 勾／1 未勾（AC#1 设计核） | `Uncommitted`／`Diff` 在同一枚 `GitView` 里加维（`git.go:178/:181`），与 181 同一文件族、没另起实现 |
| K6 状态显示 | **票 181** | 已闭 | **已结案**，7/7 勾 | 分支／工作树／repoRoot／switchBlocked 全在 `git.go`；分工见 §5 |
| K7 工具地基 | **票 163 AC#1b** | 只登记依赖，不重开 | 5 未勾／1 勾 | 该票就是"`shell.session` 未注册"的账；注册循环 5 支无 Shell 的读数是它的旁证 |
| K7 面板载体 | **票 191（甲支）** | 并格到该票 | 6 未勾／1 勾（AC#0 已裁甲） | 191 整张票就是"面板终端这一条执行面怎么裁"；甲支＝复用工具腿＋门控＋审计 |
| K8 子代理 | **票 197**（状态/流页）＋**票 188 AC#1**（产品口径已裁乙） | 并格 | 197：5 未勾／2 勾；188：5 未勾／1 勾 | `TaskRowView` 三件套是 197 落的；188 AC#1 的裁定已经落地成"进产品" |
| K9 后台任务 | **票 188 AC#2**（状态维）＋**票 164**（名册，已闭）＋**票 196**（词表，已闭） | 并格 | 188：5 未勾／1 勾；164／196 已结案 | 名册在、状态维在 `TaskRowView.Status`；剩的是 188 AC#2 的口径更新（§9-3） |
| K10 浏览器 | **票 168** | 不立新票 | 7 框全未勾；被 51 扣住 | 全池唯一认领者；AC#5 自己写明"不许新增面板可给出允许的出口" |

### 3.1 「票池里没有一张票认领它」这句票面前提，今天仍过期

尺：`grep -rlniE "内置终端|文件树|任务监控" .scratch/wisp/issues/*.md` ⇒ **8 枚**（77／182／186／188／189／190／191／196）。
与 09-28 `182-c1` 的第 36 行**同一个数**（8 枚、同一批），说明立票链一个没掉、也没有新增第 9 枚认领者。
⇒ 本票今天的价值**不在"找认领者"**，在下面这张"每一堆现在到底缺什么"的读数表。

---

## 4. AC#3 那个特别问：面板的"宿主侧读面"要不要与模型侧工具（`fs.list`／`shell.exec`）**共享实现**？

### 4.1 现量（这一问今天已经不是纯纸面问题：三枚宿主读面已经做出来了）

| 事实 | 读数 | 尺 |
|---|---|---|
| 面板调模型侧工具实现 | **0**：`grep -rn "fs\.list\|FSListDecl\|tools\.FSList" internal/panel --include=*.go \| grep -v _test`＝**0** | 正控＝`grep -rn "internal/tools" internal/panel/*.go \| grep -v _test`＝**17 行**（包在，尺能响） |
| 面板从 `tools` 包拿的只有类型 | `tools.TaskOutput` 4／`tools.TaskRoster` 2／`tools.PathCanonicalizer` 2／`tools.TaskKindRoot` 1／`tools.Decision` 1 | 同一把 `grep -oE "tools\.[A-Za-z]+" \| sort \| uniq -c` |
| 已有宿主读面怎么来的 | **各自独立、只读、不共享执行面**：181＝只读文件（`git.go:6` 逐字「reading files: is this workspace a git tree…」，全文件无 `exec.Command`）；189＝`git_readonly_exec.go`（固定查询字面，`git status`／`git diff`）；190＝`tree_reader_190.go`（`git rev-parse`＋`git ls-files`），**面板不经过 `fs.list`** | 见 §2 表 K4／K6 行的尺 |

### 4.2 两败（照旧摊清，本腿不替编排者选）

- **共享**：`fs.list` 带着 D34 的判级（L0／L2 越界）与 C26 的 allowlist 语义（`PLAN.md:2540`；`internal/tools/paths.go` 那条"空 allowlist 什么都不授权"）。
  面板一复用，**「面板显示」就走进「模型被允许读什么」的判定**；而 `AGENTS.md §1.2` 逐字禁的是"由面板侧来源的 L2『允许』"——
  共享会把那条禁令的射程从"批准"扩到"读"，**需要人工批准**。
- **不共享**：两份枚举必然漂，可能出现"面板看到的树／diff"与"模型被允许读的树"不是一棵。票 189 的 AC#2 用"同一文件族不另起实现"防的就是这个，但它防的是 **git 那一族内部**，不是跨"面板 vs 工具"。

### 4.3 本腿的读数结论（判的，不是量的）

今天三枚已落地的宿主读面**都选了"不共享实现、只共享路径/门控语义"**，而且都各自付了代价并写进文件头：
181 用文件读、189/190 各租一次 `git`（`git_readonly_exec.go:1-57` 把"为什么必须起进程"与"为什么不是工具"逐条写了）。
**唯一还活着的这一问只剩 K7**：191 甲支的原文就是"复用既有工具腿＋门控＋审计"——那是**共享门控、不共享显示通道**的第三条路，
和 181/189/190 的先例同向。⇒ **建议今后所有宿主读面按这条先例办（只读、独立实现、路径决策走 `risk.PathResolver`）**；
但按票面 AC#3，这句是**交编排者/owner 裁**的，本腿不把它当定案。

---

## 5. AC#4 与票 181 的分工（写死，且两票今天都已结案）

| 维度 | 归属 | 今天在盘上的读数 | 框态 |
|---|---|---|---|
| git **状态显示**（是不是 git 树／当前分支／detached／本地分支名册／主工作树） | **票 181** | `internal/panel/git.go:120-172`（`GitWorktree`／`GitView`）；生产装配 `Git: rt.gitView`（`run.go:921`）→ `ReadGitForWorkspace`（`git.go:298`） | **已结案**，7/7 勾 |
| git **内容审查**（未提交枚数／逐文件 diff／改了哪几枚文件） | **票 189** | 同一枚 `GitView` 的 `Uncommitted`（`:178`）＋`Diff`（`:181`）；进程面 `git_readonly_exec.go`（3 处查询调用） | **已结案**，AC#2–#6 勾；AC#1（只读设计核）未勾但设计件与裁决件都在盘 |

- **两票都不开放** ⇒ 票 182 AC#4 担心的"两票互相归口到对方已勾的格充数"今天**无处可归**：
  181 的框全勾、189 已加 `-done`，再往任何一边"归口"都会当场露馅。本程没有做、也没有建议做这种归口。
- **一处必须写死的边界**（从 181 文件头读出来的）：`git.go` 与 `git_readonly_exec.go` 是**同一族、两个文件**，
  新增的只有 `composer.git` 底下的两个**嵌套子段**（`uncommitted`／`diff`），**没有新增顶层键**——
  这是 189 为避开 `Q-51`（前端键集契约）刻意选的位置（`git.go:172-182` 的注释逐字说了）。

---

## 6. AC#5 雷区清单：面板发起的批准/写动作

**现状合规锚（本程现读，行号未漂）**：面板入向方法名册今天共 **7 枚**——
`internal/panel/bridge.go:42-45` 四枚 `panel.*`（`mode.request`／`workspace.request`／`attachment.add`／`message.send`）
＋ `:66-67` `config.get`／`config.set` ＋ `:75` `config.models.list`；
`knownComposerMethod`（`:154-158`）逐个点名同样 7 枚，**里面没有 decide／allow／commit／push／revert／stop**。
`config.set` 的写面另有一道闸：`lockedFieldFamilies`（`config_handlers.go:109`）逐字含 `plugins.`／`risk.`／`fs.`／`net.` 等族，
**插件启用开关走不到设置页的写入腿**（拒绝语逐条指回它该走的那扇门）。

| # | 雷 | demo 控件（现读） | Go 侧今天的门 | 定性 |
|---|---|---|---|---|
| M1 | **附件移除** | `rightbar.js:53`／`:74` `data-attach-del`（2 处） | 入向只有 `panel.attachment.add`，**无 remove**（`bridge.go:42-45`） | **现存写形状，待拍板** |
| M2 | **diff 接受／应用** | `rb-review.js:96` `#rb-apply-btn`「Apply 3 changes」（`:116-119` 只 toast） | 入向无 apply／accept（`knownComposerMethod` 7 枚）；189 AC#3 钉了"不许出现动作腿" | **现存写形状，待拍板** |
| M3 | **插件启用开关** | `rb-plugins.js:65`／`:75` `data-plg-toggle`（2 处，含 `on:`／`level:`） | 入向无 plugin 方法；`config.set` 把 `plugins.` 列为锁定族（`config_handlers.go:109`） | **现存写形状，待拍板**（写配置＋改档位，最硬的一枚） |
| M4 | **面板里发起执行** | `rb-terminal.js:26`／`:77` `data-terminal-add`（2 处）＋命令输入行（`:86-` 只本地回显） | **地基都不存在**：注册循环 5 支无 Shell（`run.go:614/643/670/1012/1051`） | **现存写形状，待拍板**；191 甲支是它的出口 |
| M5 | **面板「允许」** | 今日 demo **无此控件**：approval 条目只 `showScreen('approval')`（`rightbar.js:143`／`:154`／`:160`）；`grep -cE "允许\|拒绝\|批准\|approve\|allow"` 的 2 命中都是时间线/告警**散文** | 入向名册无 `decide`／`allow`（同上）；`bridge.go:63-65` 注释逐字引 `AGENTS.md §1.2` | **潜在雷（今日无控件）**：任何新增即逐字撞禁令＋`Q-49` 一族 |
| — | ~~文件树删除／移动~~ | `grep -cE "删除\|移动\|重命名\|trash\|data-del\|data-move" design/doubao/demo/rb-files.js`＝**0** | 190 AC#3 钉了只读；`tree_reader_190.go` 写调用 0 | **今日已消失**（09-28 那 6 枚里的一枚） |
| — | ~~面板 git 提交／推送~~ | `grep -nE "接受\|回滚\|提交\|推送\|commit\|push" design/doubao/demo/rb-review.js`＝**0**（只有 Apply，见 M2） | 入向名册无 commit／push；`git_readonly_exec.go` 写动词字面 0 | **今日已消失** |

**现数＝4 枚现存的写形状（M1–M4）＋1 枚潜在（M5）**；两枚 09-28 记过的雷（文件删除/移动、git 提交/推送）今天在 demo 与 Go 两侧都读不到控件/入向。
⚠ 四枚 M1–M4 **一律单列待人拍板**：`AGENTS.md §1.2` 逐字禁"由面板侧来源的 L2『允许』"，`Q-49` 一族已钉过五枚门；
**本程没有提"顺手把批准接上"**，也没有把它们并进只读交付。

---

## 7. AC#6 契约轴自证

- **`docs/PLAN.md`、`docs/specs/**`、`internal/**`、`cmd/**` 一字节未动**：本程全部命令是 `grep`／`sed -n`／`ls`／`awk`／`git log`／`git status`／`git rev-parse`，
  **零 `go` 命令、零 `gofumpt`、零写命令、零 `git add/commit`**；写面只有本文件（收尾时仍是 `??` 未跟踪，本腿按"不 git 写"没有提交它）。
- **锚漂**：起手 `232b0846c` → 收尾 `d1b1a2d22`（他人 `323-ac5` 一笔，`git diff --name-only 232b0846c..HEAD`＝`internal/config/settings_320_test.go`＋`docs/evidence/s1/323-ac5.md` 两枚），
  **没有一枚落进本程任何尺的射程**（`internal/panel`／`cmd/wisp`／`design/doubao/demo`／票面全未动）。
- 工作树里 `internal/tools/` 有几枚**别人在飞的改动**（`system_backend_windows.go`／`system_bridge.go` 等，
  与 HEAD `232b0846c` 的提交信息「23 AC#2 tools 半边」同批），**不是本程造成的**，本程没碰、也没 add 它们。
- 票面零改动：`.scratch/wisp/issues/182-*.md` 未追加、未改一字（本腿没写票面）。

---

## 8. AC#7 排程约束

- 本票**不阻塞**票 179／176／177／175 那条"超长输出读回来"的主干链；本程只读、零产码，也没有占用 `cmd/wisp` 那条共享编译面的写锁。
- **不许被拿去做"界面先做、后端后补"的理由**：本表逐堆写明"源在/不在"，源不在的堆（K1-C／K5／K7／K10）谁要动必须先有票、先有口径。
- `frontend/**` 本程**未读、未写、未引、未转述**（AC#7 逐字）；`design/**` 只读了原型 demo 的注册表与控件名（那是本仓工作树里的设计参考树，不是前端实现）。

---

## 9. 我这把尺的缺陷与未判项（具名，不许当完整表读）

1. **〔不可判〕K2 的"独立时间线要不要一行"**：十四行表里今天没有"活动流"这一行（现量 `PLAN.md:3499-3512` 共 14 行，逐行点名：
   思考中／SSE／推理／工具调用／工具调用展开／审批等待／L2 确认卡／L1 阻止／L2 原生降级／错误／Stuck／成本／已取消／注入检出）。
   我的建议（并 4/5）是**判断**，不是读数；新增一行＝改契约文字，须 owner 批。
2. **K1-C 占用条的分子/分母口径本程没量**：只量到"面板无字段"（0 命中＋正控）。167 AC#2 的原文我读了标题（"分子或分母任一未知 ⇒ 出口里必须是'未知'这一枚具名状态"），但**没复跑它的落地口径与来源**。
3. **票 188 AC#2 的前提与代码对不上（具名报回，不替裁）**：该格逐字写"给 `Record` 补状态维"，
   而 `grep -rnE "type Record struct|type Job struct|type BackgroundJob" internal --include=*.go | grep -v _test.go`＝**0**；
   正控 `type TaskRoster struct`＝**1**（尺能响）。状态维真身在 `TaskRowView.Status/StatusKnown/StatusReason`（197 落的）。
   ⇒ 该格要么改口径、要么改归属，**须编排者裁**；本程不改票。
4. **K5 的"有源"只到库层**：`memory.ListPluginStates`（`dao_misc.go:243`）与 `UpsertPluginState`（`:188`）都只有定义＋单测，
   **非测试调用者各 0**；`grep -rn "plugin/manifestslot" internal cmd --include=*.go | grep -v _test.go`＝**1 行注释**（`internal/plugin/command.go:16`），
   正控 `internal/plugin"` 非测试 import＝**8 行** ⇒ `manifestslot` **生产 import 0**。
   所以我记的是"**库里有读面、生产零装配**"，不是"运行时能取到真值"。
5. **K4／K6 的读面都租了 `git` 可执行文件**（`tree_reader_190.go:27`、`git_readonly_exec.go:349`）：这是**读数**，
   不是我对"应不应该起进程"的判断；189 AC#1 那一格（设计核）框未勾，谁要评议该去读 `189-uncommitted-count-design-a1.md`。
6. **`182-a1`/`182-c1` 转引的三处本程未复跑**：K9 的 `rightbar.js:90`、K10 的 `issues/77:364/:384`、K2 的细维出处。
   本程只复跑了主尺与各堆"源在不在"的那些尺。
7. **`PLAN.md` 行号漂移**（引用者必读）：票 145 与 145 证据件引的 `:3473-3488`／`:3481-3494` 今天**都不是十四行表**；
   现量全表＝**`PLAN.md:3499-3512`**（`## 17.5` 在 `:3492`），**行 12「成本」＝`:3510`**（旧号 `:3486`）。
8. **`ComposerState` 枚数**：今天 **12 个 json 字段**（`composer.go:271-` 区），比 `182-c1` 记的 11 多一枚（`files`，票 190 落的）。
   `Snapshot` 顶层 **6 键**（`pending/results/composer/generatedAt/instructions/tasks`），比票面立票时的 4 枚多两枚（`instructions` 200 落、`tasks` 197 落）。

---

## 10. 相对前人读数的订正（只报读数，不评功过）

| 前人的话（出处） | 今天现量 | 判 |
|---|---|---|
| `182-c1` §2：K4 文件树"无源但规格已立" | `TreeView` 7 字段＋`files` 键＋`run.go:922` 装配，190 已结案 | **已过期**（票 190 落了） |
| `182-c1` §2：K6 未提交/diff"无源但规格已立" | `GitView.Uncommitted`／`.Diff`＋`git_readonly_exec.go`，189 已结案 | **已过期**（票 189 落了） |
| `182-c1` §2：K8/K9"状态维缺" | `TaskRowView.Status` 三件套＋`tasks` 键，197 的 AC#4/#5 已勾 | **部分过期**（197 落了载体；188 AC#2 口径另说，见 §9-3） |
| `182-c1` §2：K5"规格真空、零读面" | 库里多了 `ListPluginStates`（单测在跑），但生产零装配、面板零读面 | **半过期**：读面在库层出现了，栏位仍是真空 |
| `182-c1` §4／`182-a1` §6：雷区 **6 枚** | 今天 **4 枚现存＋1 枚潜在**；文件删除/移动、git 提交/推送两枚在 demo 与 Go 两侧都读不到 | **两枚已消失**（读数见 §6） |
| 票 182 现量#1：「没有票认领它」 | 8 枚票命中（与 c1 同数同批） | **仍过期**（立票链早接住） |
| 票 182 现量#4：`git diff` 尺 0 命中 | 今天 `git_readonly_exec.go` 有 3 处查询调用、`GitView` 有 `diff` 子段 | **已过期**（189 落了） |
| 票 182 现量#5：`bridge.go:42-45` 四枚名册 | 今天 **7 枚**（+`config.get`／`config.set`／`config.models.list`） | **已扩**（313／320 一族落的） |

---

## 11. 一句话交件

**10 堆里：整堆已落 4 堆（K4 文件树、K6 审阅、K8 子代理、K9 后台任务）；部分已落 2 堆
（K1 的路径/附件两行、K3 的 pending 主体）；只差票 145 那一集载体 2 处（K1-D token、K2 活动）；
只差票 167 三个维 3 处（K1-C 占用、K3 序号、K3 停止）；无源但票已在手 1 堆（K7 终端，163 地基＋191 甲载体）；
无人认领 1 堆（K5 插件）；有票但被前置扣住 1 堆（K10 浏览器，168 等 51）。**
雷区 **4 枚现存（附件移除／diff 应用／插件开关／面板执行）＋1 枚潜在（面板允许）**，一律单列待人拍板。
本程零产码、零翻勾、零 `git` 写；写面只有这一枚文件。
