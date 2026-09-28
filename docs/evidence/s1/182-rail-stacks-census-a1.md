# 182-a1 —「任务监控」右侧竖条：堆数现数 ＋ 逐堆 Go 侧有无源 ＋ 逐堆归属票（只读普查）

- 分工：`.scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md` §A（票 182）。
- 本程**零产码**、`internal/**` 与 `docs/PLAN.md` 与 `docs/specs/**` 一字节未动、票面 AC 框**一枚没勾**。
- 写面只有两枚路径（终态闸门见 §7）：本证据件 ＋ 票 182 面 `Progress log` 追加一节。
- 预算：任务书 ≤20 枚。本程**实耗 21 枚，超 1 枚**，被哪一步吃掉具名见 §8。

---

## 0. 起手三件（共同规矩第一条，逐枚抄名册）

```
$ date
2026-09-28 17:48+0800 本机=17:48

$ git rev-parse --abbrev-ref HEAD
dev

$ git log --oneline -1        # 起手锚
38fc7c0e ledger(A387-A388) 收 33-r3（入向听众进树＋词面尺改能力尺，我四把尺复跑）＋owner 提问后按默认动作派 145-r2

$ git status --porcelain      # 起手名册（逐枚，共 63 行；终态须与此逐枚相等）
```

```text
M .gitignore
 M .scratch/wisp/probes/152/my152.py
 M .scratch/wisp/probes/161/r6/logs/flip-1.txt
 M .scratch/wisp/probes/161/r6/logs/flip-2.txt
 M .scratch/wisp/probes/161/r6/logs/flip-3.txt
 M .scratch/wisp/probes/161/r6/logs/flip-4.txt
 M .scratch/wisp/probes/161/r6/logs/flip-5.txt
 M .scratch/wisp/probes/161/r6/logs/flip-6.txt
 M .scratch/wisp/probes/161/r6/logs/flip-baseline.txt
 M .scratch/wisp/probes/161/r6/logs/flip-restored.txt
 D design/assets/base.css
 D design/assets/icons.js
 D design/assets/theme.js
 D design/assets/tokens.css
 M design/doubao/README.md
 M design/doubao/demo/app.js
 M design/doubao/demo/index.html
 M design/doubao/demo/styles.css
 D design/index.html
 D design/screens/approval.html
 D design/screens/ball.html
 D design/screens/chat.html
 D design/screens/config.html
 D design/screens/cost.html
 D design/screens/firstrun.html
 D design/screens/palette.html
 D design/screens/privacy.html
 D design/screens/security.html
 D design/screens/states.html
 D design/screens/tasks.html
 M docs/evidence/s1/152-subject-death-never-measured-r1-accept-r1.md
?? .scratch/wisp/.scratch/
?? .scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md
?? .scratch/wisp/probes/139/accept-r1/
?? .scratch/wisp/probes/152/overlay-probe1-on-samppost.json
?? .scratch/wisp/probes/156/__pycache__/
?? .scratch/wisp/probes/156/mut-156-r2/asis.log
?? .scratch/wisp/probes/156/zero156-r4-head.sh
?? .scratch/wisp/probes/156/zero156-r4-work/
?? .scratch/wisp/probes/158/r2/
?? .scratch/wisp/probes/161/r2/__pycache__/
?? .scratch/wisp/probes/161/r2/ctl/
?? .scratch/wisp/probes/161/r5/negative-control/
?? .scratch/wisp/probes/161/r6/logs/flip-7.txt
?? .scratch/wisp/probes/162/r4/
?? .scratch/wisp/probes/162/v1-baseline-gotest.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-end-bad.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-end.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-start-bad.txt
?? .scratch/wisp/probes/176/r1/logs/gate-clauses-start.txt
?? .scratch/wisp/probes/183/accept-v1/
?? .scratch/wisp/probes/185/c1/logs/d22scan-post-final.txt
?? .scratch/wisp/probes/33/r2/d22scan-final.txt
?? .scratch/wisp/probes/33/r2/gate-final.txt
?? .scratch/wisp/probes/33/r2/status-final.txt
?? .scratch/wisp/probes/33/r2/status-start.txt
?? .scratch/wisp/probes/999/
?? .zcodeignore
?? design/doubao/01-ball-states.jpg
?? design/doubao/demo/lib/
?? design/doubao/demo/rb-files.js
?? design/doubao/demo/rb-plugins.js
?? design/doubao/demo/rb-review.js
?? design/doubao/demo/rb-terminal.js
?? design/doubao/demo/rightbar.js
?? design/doubao/demo/screens/home.js
?? design/doubao/demo/screenshots/
?? design/doubao/demo/sidebar.js
?? design/old/
?? part1-state1-fixed.txt
?? part1-state1-pristine.txt
?? part1-state2-fixed.txt
?? part1-state2-pristine.txt
?? part2-nog6-fixed.txt
?? part2-nog6-pristine.txt
?? part3-stale-fixed.txt
?? part3-stale-pristine.txt
```

**起手名册枚数现数（递归尺，非 `ls | wc -l`）**：`git status --porcelain` 行尺＝**63 枚**（` M`/`M ` 15、` D` 16、`??` 32 行；枚数由上面那份逐枚清单数出，非目录级）。终态差集见 §7。

## 1. 锚：自取锚 ＝ 编排者锚，**跑到一半被推进**，具名如下

```text
$ git log --oneline -1                      # 17:5x 复跑（起手是 38fc7c0e）
971b7dad ledger(A389)＋立 Q-71 - 前端等数据的五枚票一次排开（145-r2/188-r1 写＋182a/180a/189a/190a 只读）

$ git show --name-only --format="%H %s" HEAD
971b7dad3e9392a5a52bf6b9a049be1cf5e79b58 ledger(A389)＋立 Q-71 - 前端等数据的五枚票一次排开（145-r2/188-r1 写＋182a/180a/189a/190a 只读）

.scratch/wisp/dispatches/2026-09-28-175x-wave-data-carriers-182a-180a-189a-190a-188r.md
docs/reports/pending-and-issues.md
```

- 我锚上多出来的是**编排者自己那枚 ledger(A389)**，只碰 `dispatches/**`（我的任务书，只读）与 `docs/reports/pending-and-issues.md`（台账）。**它没碰我任何一根尺的射程**（`design/doubao/demo/**` 只读、`internal/**` 只读、我的两枚写点在 §7 具名）。按共同规矩"发现 ≠ 编排者锚不算漂移"，此处只点名不复算。
- 在飞的 `145-r2`（`internal/panel/composer.go`／`pump.go`／`cmd/wisp/panel_pump.go`）**未进 HEAD**：我现读 `Snapshot` 仍是四枚（`internal/panel/composer.go:57-62`），见 §5 那条影响。

---

## 2. 表① 堆数现数（AC#1）

### 2.1 尺与对照真值

| 尺 | 命令（逐字） | 读数 |
|---|---|---|
| 主尺（递归，数注册的那一堆） | `grep -rn --include=*.js "RbPanels.register(" design/doubao/demo/` | **7** |
| 总枚数复核 | `grep -rho --include=*.js 'RbPanels.register(' design/doubao/demo/ \| wc -l` | **7** |
| 对照真值 甲（这 7 枚是否真被装载） | `grep -n "<script" design/doubao/demo/index.html` | `:173 rightbar.js`＋`:174 rb-files.js`＋`:175 rb-review.js`＋`:176 rb-terminal.js`＋`:177 rb-plugins.js` 五枚脚本在册 ⇒ `rightbar.js` 内 3 枚 ＋ 其余四文件各 1 枚＝**7**，与主尺相等 |
| 对照真值 乙（标题尺） | `grep -rn --include=*.js -E "^\s*label:" design/doubao/demo/` | **8** 命中，其中 `app.js:300` 那枚是**命令面板条目**不是竖条堆（上下文见 `app.js:9 register(id, config)` 的 palette 装配）⇒ 剔除后 **7**，与主尺相等 |

⚠ 两把对照尺（甲装载／乙标题）与主尺**三数合一＝7**。上面那把 `ls`-型尺一律没用（本程没用 `ls <目录> \| wc -l` 数任何东西）。
⚠ `design/doubao/demo/{rightbar,rb-files,rb-review,rb-terminal,rb-plugins,sidebar}.js` 与 `screens/home.js`、`screenshots/`、`lib/` 全是 **`??` 未跟踪**（`git ls-files design/doubao/demo/*.js` 只出 `app.js` 一枚）⇒ 这 7 枚堆的**当前形状在 git 之外**，本表是**现读工作树**的读数，任何人复算要在同一个工作树，不能只 `git show HEAD:`。本节只读不写、也不把它们计入任何零命中宣称（共同规矩那条）。

### 2.2 逐堆点名 ＋ "我在谁那儿看到的这一堆"（AC#1 硬要求）

| # | 堆（demo 里的逐字 label） | 注册位 `file:line` | 我在谁那儿看到的 | 内部子行（同一堆里的分行，逐枚带位） |
|---|---|---|---|---|
| K1 | 上下文 `context` | `design/doubao/demo/rightbar.js:8-10` | **demo 现注册表**（设计侧树，非本仓规格） | 工作目录 CWD `rightbar.js:20-26`；附件 `:29-33`；上下文预算 IN/OUT/CACHED＋42% 计量条 `:36-47` |
| K2 | 活动 `activity` | `design/doubao/demo/rightbar.js:83-85` | demo 现注册表 | 活动时间线 5 条样例 `:89-93`（`fs.listdir 完成`／`task.extract 运行中`／`fs.delete 被拒 L2`／`TTS 播报完成`／`审批超时自动拒绝`） |
| K3 | 审批 `approval` | `design/doubao/demo/rightbar.js:116-118` | demo 现注册表 | 待审批队列 9/10 ＋满额自动拒绝告警 `:129-138`；条目 3 条 `:121-125`；「查看全部」`:154` |
| K4 | 文件 `files` | `design/doubao/demo/rb-files.js:10` | demo 现注册表 | 工作区文件树那一堆 |
| K5 | 插件 `plugins` | `design/doubao/demo/rb-plugins.js:8` | demo 现注册表 | 分组标题「技能」`:13`／「伙伴」`:21`；条目 `:15-17`、`:23-24`（含 `level: L1/L2`、`on:` 开关位）；开关控件 `:65 data-plg-toggle` |
| K6 | 审阅 `review` | `design/doubao/demo/rb-review.js:8` | demo 现注册表 | 逐文件 diff／未提交枚数那一堆 |
| K7 | 终端 `terminal` | `design/doubao/demo/rb-terminal.js:8` | demo 现注册表 | 「新建终端」按钮 `:26 data-terminal-add` |

**表① 结论：竖条今天真有名有位的堆＝7 枚（K1–K7）。**

### 2.3 三枚"文本里有名、demo 注册表里没有位"的堆（必须分开记，否则下一位会把它们当已存在）

| # | 堆 | 谁那儿看到的（逐字出处） | demo 里有没有 |
|---|---|---|---|
| K8 | 子代理状态 | 票 182 标题与票面 AC#1 目测清单；票 188 标题逐字"subagent and background task stacks have no go-side source"；票 188 `AC#1 先答"子代理这一维到底要不要进产品"` | **无**（`RbPanels.register` 7 枚里没有） |
| K9 | 后台任务名册 | 票 188 `AC#2 后台任务状态这一维先落地`；票 164 `AC#1 名册先行` | **无独立一堆**（demo 把它混写进 K2 活动的时间线 `rightbar.js:90 task.extract 运行中`） |
| K10 | 浏览器 | `.scratch/wisp/issues/77-frontend-scaffold-reactbits-beautifului-shadcn.md:364` 逐字「`5e23d99`（右栏三标签定案：审查/终端/浏览器）」、`:384` 逐字「审查/终端/浏览器三面板、任务监控弹窗…」 | **无**（三标签定案里那枚"浏览器"在今天 7 枚注册表里已不存在） |

⇒ **这一栏全量可点名堆＝10 枚（7 在册＋3 有名无位）**。这一句是本票对 `182-c1` 那张表的**实质订正**：`182-c1` 现量也是 7，但那 7 枚（子代理／环境信息／审查 diff／文件树／内置终端／浏览器／成本 token）与今天 demo 的 7 枚**成员不同**——demo 今天多出"上下文／活动／审批／插件"四枚、少了"浏览器／子代理"两枚。**枚数巧合相同，内容不同**，所以本表才是这一波排序的依据。
⚠ owner 那张截图是**外部产品（qoder）的界面，不是本仓规格**（票 182 AC#1 原文），本表全部来源只指本仓工作树文件与票面文本，没有一条来自截图。

---

## 3. 表② 逐堆三档定性（AC#2）

档位口径：**有源**＝Go 已能取到真值（缺不缺载体另说）／**无源但规格要求**＝Go 取不到，但 `PLAN.md`／`docs/specs/**` 有人要求过（并注明要动哪一节、是不是契约）／**规格真空**＝仓内今天没人要求过，做＝新需求。

| 堆 | 档位 | 现读凭据 `file:line` | 要动哪一节／是不是契约 |
|---|---|---|---|
| K1 上下文·CWD | **有源＋已有载体** | 类型 `internal/panel/composer.go:179-193`（`WorkspaceView`，含 `Canonical`/`Reparse`/`Rewritten`/`Reason`）；进快照 `internal/panel/composer.go:207 Workspace WorkspaceView` | 不动规格（票 92 已落） |
| K1 上下文·附件 | **有源＋已有载体** | `internal/panel/composer.go:208-211`（`Attachments`/`AttachmentMIMEs`/`MaxAttachmentB`/`AttachmentReason`）；构造 `:220-225 NewComposerState` | 不动规格（票 92 已落） |
| K1 上下文·预算/token | **有源缺载体** | 源链：`internal/agent/cost.go:23 type Cost`→`:38 func (c *Cost) AddUsage(...)`→`internal/agent/guard.go:155 func (g *Guard) AddCost(...)`→`internal/agent/loop.go:976 func (l *Loop) publishUsage(...)`。面板侧零消费者：`grep -rn --include=*.go -E "publishUsage\|UsageView\|TokenUsage" internal/panel/`＝**0 命中** | D15 上下文预算（`PLAN.md:391`）／D32 那张 SLO 表——**载体＝票 145 那一集，是契约轴**（`Snapshot` 键集被前端契约钉住，见 `internal/panel/composer.go:212-215` 注释自陈） |
| K2 活动（事件流） | **有源缺载体**（部分） | 事件接口 `internal/agent/sink.go:64 Publish(Event)`、生产者 `:85 func (s *RecordSink) Publish(e Event)`；任务侧名册 `internal/tools/task.go:110-112 TaskRoster{mu, byTask map[string]TaskOutput}`——**无状态维**（本程现读该 struct 只有三行，见 §4 K9 那行）；`internal/panel/**` 里没有任何"活动／事件流"消费面（同把尺 `-E "publishUsage\|UsageView\|TokenUsage"` 0 命中的同一射程） | 载体要新枚顶层键＝**契约轴**（票 145）；事件形状本身规格已写（D10 结果呈现 `PLAN.md:242`／C21） |
| K3 审批队列 | **有源＋已有载体** | `internal/panel/composer.go:58 Pending []ApprovalCardView`（快照第一维）；面板入站方法名册 `internal/panel/bridge.go:42-45` 只有 `panel.mode.request`／`panel.workspace.request`／`panel.attachment.add`／`panel.message.send`——**没有 `decide`/`allow`**（＝今天这一堆在 Go 侧只读不批） | 不动规格；批准那一支见 §6 雷区 |
| K4 文件树 | **无源但规格要求**（宿主侧） | 模型侧**有**：`internal/tools/fs.go:309 func FSListDecl()`、`:325 BuiltinFSEntries`、注册循环 `cmd/wisp/run.go:345 tools.BuiltinFSEntries(...)`。宿主侧读面**无**（`internal/panel/` 现读文件族：`approval/assets/attachments/bridge/composer/composer_dispatch/composer_handlers/git/pump/workspace` 十支实现，**无 tree／files 一支**，尺＝`ls internal/panel/*.go` 逐名） | D34 内置工具权威表（`PLAN.md:2514` 起；`fs.list` 那一行＝L0/L2 越界）；**共享实现＝契约轴**（AC#3 那句要我裁，见 §5） |
| K5 插件 | **无源＋规格真空（面板那一堆）** | `ls internal/plugin/` 只有 `disposal.go`（＋其 test）；导出面全是卸载/释放形状：`internal/plugin/disposal.go:42 MemoryReleaser`、`:62 SLORecorder`、`:73 DisposeError`、`:89 DisposeResult`、`:98-118 Option/WithTimeout/WithRegistry/WithMemoryReleaser/WithSLORecorder`——**没有任何"列出已装插件/技能/伙伴"的读面**；`grep -rn --include=*.go -E "^func [A-Z]\|^type [A-Z]" internal/plugin/*.go \| grep -v _test.go` 全命中都在 disposal 族内 | 概念**有**规格（D3 插件运行时 `PLAN.md:108`、D19 插件信任模型 `:558`、D46 Tier-1 CLI 插件 `:2633`、票 3/19/46 那一族），但"竖条里列出插件并给开关"**仓内没人要求过** ⇒ 做＝**新需求**；开关写动作另见 §6 |
| K6 审阅 diff／未提交枚数 | **无源但规格要求** | `internal/panel/git.go` 今天只读**refs 那一族**：`:93 GitWorktree`、`:117 GitView`、`:123 Branch`、`:141 Branches`（注释逐字"refs/heads walked **RECURSIVELY**"）、`:172 ReadGit`、`:240 ReadGitForWorkspace`、`:284 locateGitDirs`、`:315 resolvePointerFile`——**不含工作树内容 diff／未提交枚数**；`grep -rn "git diff\|diff --numstat" --include=*.go internal/ cmd/ \| grep -v _test.go`＝**0 命中**（⚠ 本程未对该字面尺打正控 ⇒ 该读数按 AC#2 口径记**〔不可判〕**，写腿那一步要自带正控，见 §8 尾注） | 票 189 票面已定"同一枚读面 `internal/panel/git.go`，不另起第二份实现"（任务书 §C②）；D34/C17 白名单＝**契约轴** |
| K7 内置终端 | **无源但规格要求** | `grep -rn --include=*.go -E "shell\.session\|shell\.exec\|BrowserSession\|browser\.\|PTY\|subagent\|Subagent" internal/ cmd/ \| grep -v _test.go` ⇒ 命中**只有注释与"未接线自陈"**：`internal/config/unwired.go:63`「no shell.exec tool is registered in internal/tools, so nothing reads this flag」、`:64`、`:70`、`:76`（`lands: it lands with the shell.exec tool`）、`internal/risk/rules_shell.go:14`（R6 forced-argv 注释）、`internal/agent/spill.go:23`（1MB 上限注释）、`internal/tools/paths.go:51`（allowlist 注释）。`PTY`／`shell.session`／`subagent` **零真实现** | D34 表内逐字一行「`shell.session`｜在一条活着的会话里跑下一条命令｜**L2**｜票 163」（`PLAN.md:2514` 起那一节）＝**契约行**；票 163 已裁定该工具未注册，本票只登记依赖不重开 |
| K8 子代理状态 | **无源＋规格真空** | 同 K7 那把尺的射程内 `subagent\|Subagent` **0 真实现**（正控＝同一把尺在同一次调用里命中 `shell.exec` 六处注释，见上 ⇒ 尺本身能响） | 票 188 `AC#1` 现量它"规格真空"，任务书 §E 逐字「**不是你裁、也不是我裁**…我已摆给 owner（`A389`）」⇒ 立＝新需求，**归 owner** |
| K9 后台任务状态 | **有源缺载体缺状态维** | 名册 `internal/tools/task.go:110-112`（现读 struct 逐字只有 `mu sync.RWMutex` / `byTask map[string]TaskOutput` / `}` 三行 ⇒ **今天没有状态维**，本程亲眼量到）；输出支 `internal/tools/task.go:292-295 BuiltinTaskEntries`；注册循环 `cmd/wisp/run.go:365 tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks})`；状态维射程内零现读：`grep -rn --include=*.go -E "type Record struct\|type Job struct\|type BackgroundJob" internal/ \| grep -v _test.go`＝**0 命中**（正控＝与它同一次调用、同一模式的 `^type TaskRoster struct` 在 `internal/tools/task.go:110` **命中** ⇒ 尺能响，0 为真零） | 票 188 `AC#2 后台任务状态这一维先落地`＋票 164 `AC#1 名册先行`；对齐的是 **D43 状态机权威表 `PLAN.md:3049`＝契约**（§E 已写"不新造名字、不改表一字"） |
| K10 浏览器 | **无源＋今天连堆的位都没有** | K7 那把尺里 `browser\.`／`BrowserSession` 0 真实现；demo 7 枚注册表不含它（§2.1 主尺） | 只有 `issues/77:364/:384` 那两行前端组件清单提过（那是**组件清单不是 Go 侧数据源**，票 182 关联条原文）；`PLAN.md` 里 `web.open`/`web.fetch` 是**模型侧**工具（`PLAN.md:786`、`:2548`）⇒ 宿主侧那一堆＝规格真空，做＝新需求 |

**表② 小结（供排序）**：
- 有源（含已有载体）＝**K1 三行中的两行＋K3**；
- 有源缺载体＝**K1 token、K2 活动、K9 后台任务**（这三枚**只差票 145 那一集的载体**，是"最便宜的一批"）；
- 无源但规格要求＝**K4 文件树、K6 审阅 diff、K7 终端**（要动 D34／C17 那一族，是**契约轴**，且各已有票：190／189／163）；
- 规格真空＝**K5 插件、K8 子代理、K10 浏览器**（做＝新需求，须摆给 owner；K5 在候选票里**无人认领**）。

---

## 4. 表③ 逐堆指认归属票（AC#3＝本票主要产出）

候选池（任务书 §A 点名）＝票 145／181／188／189／190／174／164。七枚的 AC 标题本程**现读**（尺 `grep -h -m3 -oE "AC#[0-9]+ [^*]{0,60}" .scratch/wisp/issues/<n>-*.md`），且 `-done` 计数尺 `ls .scratch/wisp/issues/ | grep -c "^<n>-.*-done"` 对七枚**全为 0** ⇒ 七枚今天**都没结案**，并格不会踩到"已勾格"。

| 堆 | 归属 | 具名并格（并进谁的哪一格） | 这一维该进快照的哪一格 |
|---|---|---|---|
| K1 上下文·CWD | **已闭，无需并格** | 票 92 落的（现读 `internal/panel/composer.go:207`） | 已在 `Composer.workspace`（`:207`） |
| K1 上下文·附件 | **已闭** | 票 92 | 已在 `Composer.attachments/acceptedAttachmentMimes/maxAttachmentBytes/attachmentError`（`:208-211`） |
| K1 上下文·预算/token | **票 145 扩快照** | 并 145 `AC#1 先把"哪一行缺哪个字段"落成一张对照表` 的**成本行＝行 12**（现读出处 `docs/evidence/s1/145-snapshot-field-census-r1.md:378` 节标题逐字「### 2.13 行 12 · 成本（`:3486`，C23）」）；⚠ 145 面 `AC#2/AC#6 维持未勾`（票面自陈）⇒ 并格合法但**别当已勾** | 新增顶层键（`Cost`／`Budget`）——**键集是契约**（`internal/panel/composer.go:212-215` 注释逐字：新增键要前端契约同批、`Q-51`） |
| K2 活动（事件流） | **票 145（载体）＋票 188（内容）** | 载体并进 145 `AC#1` 那张表要**新增一行**（十四行表里今天没有"活动流"，票 182 现量#2 已证）；"运行中／被拒 L2／超时"这些**状态词**归 188 `AC#2`（逐字对齐 D43） | 新枚顶层键 `Feed`（＝`Snapshot.Pending`/`Results`/`Composer`/`GeneratedAt` 之外的第五维，现读 `internal/panel/composer.go:57-62` 只有四枚） |
| K3 审批队列 | **票 145 射程外，已存在** | 谁都不并：`Pending` 已在快照（`:58`）。145 那张十四行表要**别把它重复列一行**（防双记账） | 已在 `Snapshot.pending` |
| K4 文件树 | **票 190** | 并 190 `AC#1 只读设计核`＋`AC#2 只读落地＋常驻判据`（现读标题逐字）；`AC#3 只读边界`是本票那一堆的雷区出口 | 新枚顶层键 `Files`；**不许走 `Composer`**（`Composer` 键集已被前端契约钉死，`composer.go:212-215`） |
| K5 插件 | **候选池里没有归属票 ⇒ 须立新票** | 145/181/188/189/190/174/164 的 AC 标题本程逐枚读过，**没有一枚射程覆盖"列出插件/技能/伙伴"**；最近的三枚概念票（D3/D19/D46 的落地票）都不在候选池，且票 186 `:48` 逐字不认领任务监控那一栏（`182-c1` 已核，本程不复算） | 新枚顶层键 `Plugins`——但**先要 owner 拍"这一堆要不要进产品"**（＝规格真空，做＝新需求） |
| K6 审阅 diff／未提交枚数 | **票 189（内容）＋票 145（载体）** | 并 189 `AC#1 只读设计核`／`AC#2 只读落地`／`AC#3 只读边界（命门）`；读面**同一枚 `internal/panel/git.go`**（任务书 §C②，与票 181 不另起第二份）；**与票 181 的分工见 §5.2（AC#4）** | 挂在已有的 `Composer.git`（现读 `internal/panel/composer.go:216 Git GitView`＝票 181 已落）里**加维**：`GitView`（`internal/panel/git.go:117`）现无 `Uncommitted`/`Files` 字段 ⇒ 189 落的是**同一 struct 的新增字段**，不是新顶层键（这样才不动键集契约） |
| K7 内置终端 | **票 163（工具地基）＋立新票（面板载体）** | 地基只登记依赖：`internal/config/unwired.go:63-76` 自陈"shell.exec 未注册"、票 163 已裁 `shell.session` 未注册 ⇒ 本票**不重开 163**；"面板里那终端那一堆的宿主侧读面"**候选池无人认领** ⇒ 与 K5 同判：立新票，且必须带 §6 那枚雷区 | 新枚顶层键 `Terminals`；⚠ 它天然带**执行**，不许并进只读那一堆交付 |
| K8 子代理状态 | **票 188** | 并 188 `AC#1 先答"子代理这一维到底要不要进产品"`——那一格**已由编排者摆给 owner（`A389`／`Q-71`）**，任务书 §E 逐字"不是你裁、也不是我裁"⇒ 本票只把"堆存在"记进表，**不替它做决定** | 若 owner 批：进 K2 的 `Feed`（同一维，别开两枚键）；若不批：竖条那一格永远显示不出来，前端要按"无此堆"实现 |
| K9 后台任务状态 | **票 188 ＋ 票 164 ＋ 票 174** | 状态维＝188 `AC#2`（现读 struct `internal/tools/task.go:110-112` 三行，确实**缺字段不缺名册**，与 188 `AC#2` 票面"名册已在、缺字段"**一致**）；名册与输出支＝164 `AC#1 名册先行`（注册证据 `cmd/wisp/run.go:365`）；输出指向模型默认读不到的文件＝174 `AC#1 真读数` | 进 K2 `Feed` 里每条的 `status` 字段，枚举**逐字取 D43 权威表 `PLAN.md:3049`**（§E 口径：不新造名字、不改表一字） |
| K10 浏览器 | **候选池无归属 ⇒ 立新票或明确不做** | 只有 `issues/77:364/:384` 前端组件清单提过；Go 侧零源（K7 那把尺）；`web.open`／`web.fetch`（`PLAN.md:786`、`:2548`）是**模型侧**，不是宿主读面 ⇒ 做＝新需求，**且今天 demo 已不给它留位**（§2.1 主尺 7 枚无 browser）⇒ 建议按"文本残留"处理、不排这一波 | 不占快照格 |

**表③ 小结（这一波排序的直接输入）**：
1. **只差载体的一批（最便宜，全归票 145 那一集）**：K1 token、K2 活动、K9 状态维 —— 三枚都能在 `145-r2` 正在长的那一集里加维，前提是先解 `composer.go:212-215` 自陈的**键集契约**（`Q-51`）。
2. **要动契约轴、但票已在手**：K4→190（本波 §D）、K6→189（本波 §C）。
3. **地基不在、票已裁不重开**：K7→163（工具未注册）、K9 名册→164、K9 输出可读性→174。
4. **候选池无人认领、必须立新票或摆给 owner**：K5 插件（开关＝写动作）、K10 浏览器（demo 已不给位）、K8 子代理（已在 `A389`／`Q-71`）。

---

## 5. AC#3 那句要我裁的问题 ＋ AC#4 与票 181 的分工

### 5.1 宿主侧读面要不要与模型侧工具共享实现？——**本程不选，两难摊清交编排者／owner 裁**
- **共享的代价**：`fs.list` 今天带着 D34 的风险分级（L0／L2 越界，`internal/tools/fs.go:309 FSListDecl`＋`internal/tools/paths.go:51` 空 allowlist 什么都不授权那条注释），一旦面板通道复用同一支，**面板的请求就走进模型侧的越界判定**；而 `AGENTS.md §1.2` 逐字禁"由面板侧来源的 L2『允许』"（`Q-49` 那一族已钉过五枚门）。共享把这条禁令的**射程从"批准"扩到"读"**，需要人工批准。
- **不共享的代价**：两份文件遍历必然漂（票 189 的任务书 §C② 对 git 读面写的正是同一条理由），越界规则也可能漂，出现"面板看到的树与模型被允许读的树不是一棵"。
- **既有先例**：票 181 选的是**第三条路**——宿主侧另起一支**只读、不执行外部命令**的实现（`internal/panel/git.go:284 locateGitDirs`／`:315 resolvePointerFile`／`:358 parseGitdirFile` 全是读文件，`grep` 该文件无 `exec.Command` 命中，且 `:6` 注释逐字"reading files and nothing else"），载体接在 `cmd/wisp/panel_pump.go:109 panel.ReadGitForWorkspace(rt.workspaceView())`。⇒ **建议 K4 文件树照 K6 的先例办**（另起只读支、不共享），但**这句只是指认先例，不是裁定**。
- ⚠ 无论走哪一支，路径解析必须留在 `risk.PathResolver` 射程内（`AGENTS.md §1.2` 逐字禁在它之外用 `filepath.Clean|Abs` 做文件系统决策；`internal/tools/task.go:105-109` 注释是同一口径的先例："C26 stays out of it, and so does the ban on filepath decisions outside risk.PathResolver"）。

### 5.2 AC#4 与票 181 的分工（写死）
- git **状态显示**（是不是 git 树／当前分支／detached／本地分支名册／主工作树）＝**票 181，今天已落地**：`internal/panel/git.go:172 ReadGit`、`:141 Branches`、`:93-101 GitWorktree`、载体 `internal/panel/composer.go:216 Git GitView`＋生产者 `cmd/wisp/panel_pump.go:109`。
- git **内容审查**（未提交枚数／逐文件 diff）＝**票 182 的 K6 那一堆 → 归票 189**。
- ⚠ 两票**不许互相"归口到对方已勾的格"充数**：票 181／189 的 `-done` 计数本程现读**都是 0**（§4 那把尺），谁也没结案；且票 182 面 AC 框一枚没勾。
- ⚠ 一处**必须提醒写腿**的形状：`internal/panel/git.go:113` 注释逐字「remoteBranches is deliberately NOT a field」——189 若要加"未提交枚数"，是在**同一枚 struct 加新字段**，不动键集契约那一枚 `Snapshot`（见 §4 K6 行），这样才避开 `Q-51`。

---

## 6. AC#5 雷区清单（逐堆：这一堆里有没有**由面板发起**的批准／写动作）

| 堆 | 有没有面板侧写动作 | 现读凭据 | 处置 |
|---|---|---|---|
| K1 上下文 | ⚠ **有两枚** | demo `rightbar.js:16 data-ctx-add`「添加」按钮、`:74 data-attach-del`「移除附件」＋`:53-58 onMount` 真删行 | 附件那一族**票 92 已有合法入口**（`internal/panel/bridge.go:44 MethodAttachmentAdd = "panel.attachment.add"`＝受控入站方法，非批准）⇒ 「添加／移除」归票 92 射程，**不单列新雷**；但"移除"在 Go 侧有没有对应方法名册里**没有**（`bridge.go:42-45` 四枚里无 `attachment.remove`）⇒ **待人拍板** |
| K2 活动 | 无 | `rightbar.js:112 onMount(app, el) { /* 纯展示 */ }` 逐字"纯展示" | 合规，可进只读那一堆 |
| K3 审批队列 | 今天无，**天然最容易越界** | demo 条目点击只跳屏 `rightbar.js:159-162 showScreen('approval')`；Go 侧 `internal/panel/bridge.go:42-45` **没有 `decide`/`allow`**（＝面板批准那条线今天不存在） | ⚠ 任何"在面板里点批准"＝**逐字撞** `AGENTS.md §1.2`"由面板侧来源的 L2『允许』"禁令＋`Q-49` 那一族 ⇒ **单列待人拍板，不并进只读交付** |
| K4 文件树 | ⚠ 树本身只读，**但删除／移动是雷** | 本程未在该 demo 文件里读到删除控件（未逐枚复算＝本程没读完，写腿要现读 `rb-files.js` 的 `data-*` 名册） | ⚠ 沿用 `182-c1` 表 §12 那一格（面板文件树删除／移动）：单列待人拍板；票 190 `AC#3 只读边界` 是它的出口 |
| K5 插件 | ⚠ **有，而且是最硬的一枚** | demo `rb-plugins.js:65 data-plg-toggle title="切换启用"` ＋条目自带 `on:`（`:15-17`、`:23-24`）与 `level: L1/L2` | 面板侧切换插件启用＝**写配置＋改风险档位**，与"档位显示与发起请求"那条 R20 口径相邻 ⇒ **单列待人拍板**，且这一堆整体是规格真空（§3 K5） |
| K6 审阅 diff | ⚠ **有两枚**（接受／回滚） | 本程未逐枚复算 `rb-review.js`；`182-c1` 已记 4 处雷（含"diff 接受／回滚"），该格本程**不复算、只转记出处** | ⚠ 接受／回滚＝写动作＋面板来源 ⇒ 单列待人拍板；票 189 `AC#3 只读边界` 是它的出口 |
| K7 终端 | ⚠ **有**：新建终端＝起进程 | demo `rb-terminal.js:26 data-terminal-add title="新建终端"`；Go 侧地基零（`internal/config/unwired.go:63-76` 自陈未注册） | ⚠ 面板内置终端＝**在面板里发起执行**（最硬的一枚，且地基都不存在）⇒ 单列待人拍板，不许并进只读那一堆 |
| K8 子代理／K9 后台任务 | 无（K9 只有"看"） | 名册 struct 无写面（`internal/tools/task.go:110-112`） | 可进只读那一堆 |
| K10 浏览器 | 无位（demo 已不给它位置） | §2.1 主尺 7 枚无 browser | 不排这一波 |

**雷区枚数现数＝6**（K1 附件 remove 1＋K3 批准 1＋K4 删除/移动 1＋K5 插件开关 1＋K6 接受/回滚 1＋K7 起终端 1；尺＝上表逐堆具名，出处每格带位；K4/K6 两格标了"本程未复算"，转记 `182-c1`，**不得当作本程读数**）。⚠ 这六枚**一律单列为待人拍板项**，本程**没有**提"顺手把批准接上"。

---

## 7. 只读自证与终态闸门（不用区间 diff）

- 自证口径：本程**没跑过任何写命令**，除 §0/§1 的 `git log`／`git show --name-only` 外全是 `grep`/`ls`/`sed -n`（只读）；**没有用区间 diff 自证**（共同规矩那条）。
- 本程写面逐枚具名（两枚）：
  1. `docs/evidence/s1/182-rail-stacks-census-a1.md`（本文件，新建）
  2. `.scratch/wisp/issues/182-*.md`（**只追加** `Progress log` 一节，票面原文一字未改、AC 框一枚未勾）
- `frontend/**`：**没读、没引、没写**（票 182 AC#7 逐字"不读不写不引不转述"）；只现数过它的**存在与形状**（`ls frontend/` 15 行，含 `dist`/`src`/`node_modules`），该枚数**不进入任何零命中宣称**。
- `design/**`：**只读不写**；§2 那把主尺扫的是 `design/doubao/demo/**`（读），其命中**不用于任何"Go 侧无源"宣称**——所有"无源"宣称的尺射程一律是 `internal/ cmd/` ＋ `--include=*.go`（逐把见 §3 表内命令）。
- 终态闸门：`git status --porcelain` 须与 §0 起手名册**逐枚差集为空**（差集＝本文件＋那枚票面两处，加进 commit 后从名册消失；验证原文与结果见本程交件回报，共同规矩"不是必须为空"）。

## 8. 预算、未裁完与停手上报项

- **预算 21 枚／上限 20 枚，超 1 枚**，被哪一步吃掉具名：第 17 枚那次调用把四段尺串在一条 `&&` 链上，链中"memory Record + jobs"那把尺**零命中即退出码 1**，导致后半段（票面 AC 标题、`145` 证据件行 12、`77` 三标签）**整段没跑**，只能补第 18 枚重跑。⇒ 教训已记：**零命中尺不许串 `&&`**。
- **本程未裁完的格子（不许当完整表交，具名）**：
  - K4／K6 两堆的**面板侧控件名册**（`rb-files.js`、`rb-review.js` 里 `data-*` 逐枚）没读完 ⇒ 雷区那两格转记 `182-c1` 并标了出处；
  - `grep "git diff|diff --numstat"` 与 `ls internal/panel/*.go` 那两把**字面尺未打正控** ⇒ 相应 0 命中按 AC#2 口径记〔**不可判**〕，写腿（§C 那枚）要自带"种 X 必响"的正控；
  - 票 145 那张十四行表的**逐行行号**只核到"行 12＝成本"这一枚（`docs/evidence/s1/145-snapshot-field-census-r1.md:378`），其余行未复算。
- **停手具名上报（不替编排者改规格）**：
  1. **前提冲突**：票 182 `AC#7` 逐字 `frontend/**` "不读不写不引不转述"，任务书共同规矩逐字"读可以，写不行"。本程**按票面从严**（`frontend/**` 一字未读）。若编排者认为堆数应当从 `frontend/src` 现数（那才是产品真身，demo 是设计参考树），请**具名改这一格**再复跑本尺——本表 §2 因此是**从 demo 注册表**数的，这条射程差异是这一波排序里最该先定的一件事。
  2. **§E（188-r1）的前提可能对不上代码**：任务书 §E 写"`Record` 补状态维，写面 `internal/memory/**`"，但本程现读 `grep -rn --include=*.go -E "type Record struct|type Job struct|type BackgroundJob" internal/ | grep -v _test.go`＝**0 命中**（同一把尺的正控 `^type TaskRoster struct` 在 `internal/tools/task.go:110` 命中，尺能响）。⇒ 状态维的真身在 `TaskRoster`/`TaskOutput`（`internal/tools/task.go`）而**不在 `internal/memory` 的 `Record`**；写腿起手现读若找不到 `Record`，那不是它写歪，是任务书那一格指错了模块。**这条要编排者裁，本程不动 §E 也不动票 188。**
  3. **K5 插件那一堆在候选池里没有归属票**（§4 表③），本程**不替它造落地票**（票面"不许我凭空替它造落地票"那句仍然生效），只点名"须立新票或摆给 owner"。
