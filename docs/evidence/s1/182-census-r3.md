# 182 — 右上角「任务监控」栏 10 堆的 r3 只读重跑（起手锚 `edab04ab6`＝派单锚，收尾同一锚；实际 HEAD `ac2c1d903`）

> 票 182 **AC#2／AC#3** 的 r3 可判证据（本腿**不翻任何框**，勾归编排者；**不触台账**）。
> **只读腿**：`internal/**`、`cmd/**`、`docs/PLAN.md`、`docs/specs/**`、`frontend/**` **零字节写入**（本程连 `frontend/**` 的读都没发生——r2 的那一次越界自报不重犯）；`design/**` 只发生 grep 读（demo 注册表口径是编排者 09-28 钉死的本票主尺）。
> 本件是 `docs/evidence/s1/182-census-r2.md`（下称 r2，锚 `5183f9715`→`a42d6fccc`）的**重跑与订正**。r2 的读数在今天的树上**全部复现**，唯一订正见 §6（一处计数伪差，非树移动）。
> **口径**：堆数沿用编排者 09-28 裁定＝原型 `design/doubao/demo/**` 注册表口径，**不是产品真身**；owner 截图是外部产品（qoder）的界面，不是本仓规格。
> **正控纪律**：凡 0 命中，同一次调用里打了已知存在的同形尺，读数记括号里；没打正控的不记数。
> **树移动证据**：`git diff --name-only a42d6fccc..HEAD -- internal/ cmd/ design/`＝**恰好 1 枚**（`internal/tools/pointer_185_reread_chain_test.go`，票 185 的新测试件，不碰本题任何源文件）；`go.mod`/`go.sum` diff＝**0**。⇒ 本题射程内的生产代码自 r2 锚以来**字节不动**，r2 的门禁读数（build rc=0、panel 两枚已知红、d22scan rc=0）**按派单"none needed"不重跑、按"树未动"承用**。

---

## 1. AC#1 复核：堆数 = 10（7 枚有位 + 3 枚有名无位），成员与 r2 完全一致

| 尺 | 原文 | 今天读数 |
|---|---|---|
| 主尺 | `grep -rn --include=*.js "RbPanels.register(" design/doubao/demo/` | **7 行**，`file:line` 与 r2 逐枚同位（rb-files:8／rb-plugins:6／rb-review:6／rb-terminal:6／rightbar:8/:83/:116） |
| 对照尺 | `grep -rn --include=*.js -E "^\s*label:" design/doubao/demo/` | **8 行**，剔 `app.js:300`（命令面板条目）＝**7**，与主尺相等 |
| 有名无位·子代理 | `grep -rn "子代理" design/doubao/demo --include=*.js` | **0**（K8 无位仍成立） |
| 浏览器有位吗 | `ls design/doubao/demo/ \| grep -iE "browser\|agent\|task"` | **空**（K10 无位仍成立） |

**7 枚有位**＝K1 上下文／K2 活动／K3 审批／K4 文件／K5 插件／K6 审阅／K7 终端（label 逐字与 r2 同）；
**3 枚有名无位**＝K8 子代理（票 188 AC#1 已裁「乙：进」＋票 197）、K9 后台任务名册（票 188 AC#2＋164/196）、K10 浏览器（`issues/77:364`/`:384` 仓内文本）。
⇒ **可点名总堆数 = 10，同数同成员**，三程（r1/r2/r3）一致。

---

## 2. AC#2 逐堆三档定性（每堆尺与命中数为今天现跑）

三档定义（票面原文）：**有源**＝Go 里已能取到真值（缺不缺载体另说）／**无源但规格要求**＝Go 取不到但 `PLAN.md`／specs 有人要求过／**规格真空**＝今天没人要求过。

| 堆 | 三档 | 现读尺（原文）→ 今天命中数（括号＝正控） | 载体 | 归属票的框态（现读） |
|---|---|---|---|---|
| **K1-A 路径（CWD）** | **有源＋载体已在** | `grep -n "Workspace .*json:\"workspace\"" internal/panel/composer.go`＝**1**（`:273`）；`awk '/^type ComposerState struct/,/^}/' … \| grep -c "json:"`＝**12** | 快照 `composer.workspace`；装配 `run.go:930` | 票 92 `-done`（5 勾／1 未勾，见 §6-①） |
| **K1-B 附件** | **有源＋载体已在** | `composer.go:274-277` 四枚（attachments／acceptedAttachmentMimes／maxAttachmentBytes／attachmentError） | 快照 `composer.*` | 票 92 同上 |
| **K1-C 占用条** | **两半源都在、口径未定、载体缺** | 面板侧 `grep -rnE "usedTokens\|UsedTokens\|contextUsed\|ContextUsed" internal/panel --include=*.go \| grep -v _test`＝**0**（正控同包 `json:`＝**191** 行）；`grep -rnE "Budget" internal/panel …`＝**0**；分母候选＝`budgets.go:79 TokenBudget`／`:50 Budgets`／`:54 ContextWindow`；分子＝`cost.go:52 AddUsage` | 无：`PumpSources`（`pump.go:138-219`）**11 个 reader**（Verdicts/Mode/Workspace/Git/Files/Model/Credential/Results/Instructions/Tasks/L1Windows）里没有占用 reader；生产装配 `run.go:927-959` 逐枚对上 | 票 **167 AC#2**（1 勾／6 未勾）：口径未定 ⇒ 不许先落字段 |
| **K1-D token／成本** | **有源缺载体** | 源链 `cost.go:52 AddUsage` → `guard.go:148 AddUsage` → `loop.go:536`/`:540` → `:1115 publishUsage`；面板消费 `grep -rniE "\bUsage\b\|\bCostView\b\|tokensIn\|totalTokens" internal/panel …`＝**0**（正控同上 191 行） | 无：`Snapshot` 6 键里没有 cost；`PumpSources` 无 Cost reader | 票 **145 行 12「成本（C23）」**＝`PLAN.md:3510`（表 `:3499-3512`）；145 **4 勾／3 未勾**（AC#2／AC#2b／AC#6）；单位口径 ⚠ 仍在 `cost.go:17` 自陈「open question」 |
| **K2 活动** | **有源缺载体（细维）** | 源＝`sink.go:51 type Event struct`（12 维）；面板 `grep -rnE "json:\"feed\"\|ActivityView" internal/panel …`＝**0**（正控 `json:"rows"`＝**1**，`subagent_roster_197.go:158`） | `Snapshot.results` 只有 `ResultChunk` 三枚，无工具名/状态/时刻 | 内容维规格＝`PLAN.md:3502` 行 4 chip ⇒ 并 **145 行 4/5** 载体扩张；独立时间线＝改 `PLAN.md`（须人工批） |
| **K3 审批（主体）** | **有源＋载体已在** | `ApprovalCardView` **10** 个 json 字段（`approval.go`）；真源 `pending_read.go:58 Position`／`:127` 赋值 | 快照 `pending`（`Snapshot` 第一维） | 票 **145 行 6/7**（已在，别重复列） |
| **K3 序号／停止** | 序号：源在载体缺；停止：源未建 | `grep -rn "Position" internal/panel …`＝**0**（正控同文件 `CorrelationID`＝**3**） | 无 | 票 **167 AC#3**（序号）／**AC#4**（停止），1 勾／6 未勾 |
| **K4 文件树** | **有源＋载体已在** | `TreeView` **7** 个 json（`tree_types_190.go:50-58`）；`composer.go:285 Files *TreeView`；装配 `run.go:932`；只读 `grep -nE "os\.Remove\|os\.Rename\|os\.WriteFile\|os\.Create\|RemoveAll" internal/panel/tree_reader_190.go`＝**0**（正控 `os.Open\|os.Lstat`＝**2**） | 快照 `composer.files`（omitempty） | 票 **190 已结案** 6/6；⚠ 跨树红（前端声明欠 `files`，`Q-51` 族）非 Go 源缺 |
| **K5 插件（读）** | **库里有读面／生产零装配／栏位规格真空** | `dao_misc.go:243 ListPluginStates`；非测试调用者 **0**（正控同文件 `func (s *Store)`＝**12**）；面板 `grep -rni "plugin" internal/panel …`＝**3 行**（1 锁定族名单＋2 散文，零读面）；`manifestslot` 生产 import＝**0**（正控 `internal/plugin"` 非测试 import＝**9** 行） | 无 | 票 **407**（open，自陈 Blocked by 前端 face 归外部 agent） |
| **K6 审阅** | **有源＋载体已在** | `GitView` **12** 个 json（`git.go:144-181`，`Uncommitted`＝`:178`、`Diff`＝`:181`）；`runGitReadOnly(`＝**4** 次（定义 `:349`＋调用 3）；写动词字面 **0**（正控查询字面 `status/diff/rev-parse/ls-files`＝**5**） | 快照 `composer.git.uncommitted`／`.diff` | 票 **189 已 `-done`**（5 勾／1 未勾）；状态显示半＝**181 `-done`** 7/7 |
| **K7 终端** | **无源但规格要求** | 注册循环 5 支＝`run.go:624` FS／`:653` Task／`:680` System／`:1038` Subagent／`:1077` Plan，**无 Shell 一支**；`grep -rn "shell\.session\|shell\.exec" internal cmd --include=*.go \| grep -v _test.go`＝**6 行全注释**（unwired.go:63/64/70/76、rules_shell.go:14、spill.go:26） | 无 | 契约行＝D34 `PLAN.md:2575`「`shell.session`｜…｜L2｜票 163」；地基＝**票 163 AC#1b**（1 勾／5 未勾）；载体＝**票 191**（1 勾／6 未勾，AC#0 已裁**甲**） |
| **K8 子代理** | **有源＋载体已在** | `TaskRowView` **12** 个 json（`subagent_roster_197.go:111-142`，`Status/StatusKnown/StatusReason`）；装配 `run.go:953 Tasks: rt.taskRosterState` | 快照 `tasks` | 票 **197**（2 勾／5 未勾）＋票 188 AC#1（乙） |
| **K9 后台任务** | **有源＋载体已在** | `TaskRosterSection` **7** 个 json（`:155-183`：rows/inFlightSlots/poolCap/三枚流截断读数） | 同上 `tasks` | 票 **188 AC#2**（1 勾／5 未勾；⚠ 该格前提「给 `Record` 补状态维」与代码对不上——`type Record struct` 尺＝**0**，正控 `^type TaskRoster struct`＝**1**（`internal/tools/task.go:176`），状态维真身在 `TaskRowView`，须编排者裁口径）＋164/196（皆 `-done`） |
| **K10 浏览器** | **无源＋规格真空（有票，被前置扣住）** | `grep -rli "browser" internal/panel …`＝**0**（正控 `grep -rli "WebView" cmd/wisp …`＝**14 文件**） | 无 | **票 192**（1 勾／5 未勾，AC#0 已定案「零 key 抓结果页」）＋**票 168**（0 勾／7 未勾，能力侧；50/51 已 `-done`，真卡 C17/C24 未定案） |

### 2.1 小结（按「只差什么」排，与 r2 同形）

1. **只差载体（票 145 那一集）**：K1-D token／成本、K2 活动。
2. **只差口径或字段维（票 167）**：K1-C 占用、K3 序号、K3 停止。
3. **已闭、只须记账**：K1-A/B（92）、K3 主体、K4（190）、K6（189＋181）、K8/K9（197＋188＋164/196）。
4. **无源、票已在手但不许本票重开**：K7（163 地基＋191 甲载体）。
5. **认领已立但仍开**：K5 读＝407（Blocked by 前端 face）；K10＝192（＋168 能力侧）。

---

## 3. AC#3 逐堆归属票（本票主要产出）

> 指认前逐张票现读了框态（读数见 §2 末列）。**本程没有替任何票翻勾，也没有改任何既有票的一个字节。**

| 堆 | 归属（具名＋格） | 并格还是立票 | 一格一句为什么是它 |
|---|---|---|---|
| K1-A／K1-B | **票 92**（已 `-done`） | 已闭 | `workspace`＋附件四枚都是票 92 落的 |
| K1-C | **票 167 AC#2** | 并格 | 该格原文就是「分子或分母任一未知 ⇒ 出口必须是『未知』具名状态」 |
| K1-D | **票 145 行 12**（`PLAN.md:3510`） | 并格（145 扩快照） | 行 12 逐字「成本（C23）」；源链现读可得、面板零消费者 |
| K2 | **票 145 行 4/5**（载体）＋**票 196**（状态词表，`-done`） | 并格；**不另立票** | 行 4 chip 已把「运行中/成功/失败/被拒＋L0/L1/L2」定义完 |
| K2 可选独立时间线 | **须 owner 一句话** | 若做＝改 `PLAN.md` 十四行表 | 加行＝改契约文字，只摆 owner 批 |
| K3 主体 | **票 145 行 6/7** | 不并（已在） | `pending` 已是快照第一维，防双记账 |
| K3 序号／停止 | **票 167 AC#3／AC#4** | 并格 | `LiveApprovals.Position` 真源在 `pending_read.go:58/:127` |
| K4 | **票 190**（已 `-done`，6/6） | 已闭 | `TreeView`＋`files` 键＋生产装配全在，写调用 0 |
| K5 读 | **票 407**（2026-11-05 立票，来路＝本普查） | 已立 | 全池唯一认领者；Blocked by 前端 face ⇒ 本编队可写面只剩「Go 侧 bridge 读口若缺另立子格」 |
| K5 开关（写） | **无票；不进本票交付** | 单列待人拍板 | 写配置＋改档位＝`AGENTS §1.2` 禁令邻域（§5 M3） |
| K6 内容审查 | **票 189**（`-done`，5 勾/1 未勾） | 已闭 | `Uncommitted`/`Diff` 加在同一枚 `GitView` 嵌套子段，没另起实现 |
| K6 状态显示 | **票 181**（`-done`，7/7） | 已闭 | 分支/worktree/repoRoot 在 `git.go:144-171` |
| K7 工具地基 | **票 163 AC#1b** | 只登记依赖 | 该票就是「`shell.session` 未注册」的账；注册循环 5 支无 Shell 是旁证 |
| K7 面板载体 | **票 191（甲支）** | 并格 | 191 整张票就是「面板终端执行面怎么裁」，AC#0 已裁甲（票面 `:3`/`:15` 现读） |
| K8 | **票 197**（载体）＋**票 188 AC#1**（乙：进） | 并格 | `TaskRowView` 三件套是 197 落的 |
| K9 | **票 188 AC#2**（口径待裁）＋**票 164**／**票 196**（皆 `-done`） | 并格 | 名册在、状态维在 `TaskRowView.Status`；剩 188 AC#2 的口径更新 |
| K10 面板 | **票 192**（AC#0 已定案「零 key 抓结果页」，票面 `:3`/`:12` 现读） | 不立新票 | 192 AC#1–#5 写的就是预览/注入面/只读边界/面板载体 |
| K10 能力/插件 | **票 168** | 不立新票 | 168 管「逛网页」能力与插件门控，自陈不做给人看的界面 |

**级联立票计数 = 0**：10 堆全部有具名归属，**本程没有新建任何 issue 文件**。不需要立票的两种候选情形都不成立——(a) 无主堆：无（r2 时 K5 已由 407 接住）；(b) 雷区写形状（M1–M4/M5）：按 r1/r2/c1 三程一致的纪律，这些是**待人拍板项**，立票＝把「面板发起的 L2 批准/写」变成可派工单，恰好越过 `AGENTS §1.2` 需人工批准的那条线，⇒ 继续单列（§5），不立票。
另两处**不替裁、具名报回**（沿用 r2）：①票 192 与票 168 的票面互相把对方那一半写成「没票」，口径应为「K10 面板＝192、K10 能力侧＝168」，须编排者收敛票面文字；②票 188 AC#2 的「`Record`」口径过期（§2 表 K9 行），须编排者裁。

### 3.1 AC#3 特别问：面板的「宿主侧读面」要不要与模型侧工具（`fs.list`/`shell.exec`）共享实现？

权威尺现跑（环境按派单：`export PATH="$PWD/build:…"`；`unset GOTMPDIR TMPDIR TEMP TMP`；未设 CGO）：

| 事实 | 尺 | 今天读数 |
|---|---|---|
| 面板生产代码 import `internal/tools` | `go list -f "{{join .Imports \" \"}}" ./internal/panel \| tr " " "\n" \| grep -c "internal/tools"` | **0** |
| 同一把尺正控 | 同上 `grep -c "internal/risk"` | **1**（尺能响） |
| `go.mod` 自 r2 锚以来 | `git diff --name-only a42d6fccc..HEAD -- go.mod go.sum \| wc -l` | **0** |

⇒ r2 的订正维持：面板与模型侧工具包**生产依赖为 0**（唯一 import 在测试件），模型侧类型经 `cmd/wisp` 装配点以值/字符串递进来。今天的既成事实＝**不共享实现、只共享路径/门控语义**（181 纯文件读、189/190 各租一次 `git` 只读进程、190 文件树不经 `fs.list`；`git.go` 全文 `exec.Command`＝**0**）。共享会把 D34 判级与 C26 allowlist 语义渗进面板读通道，撞 `AGENTS §1.2`「由面板侧来源的 L2『允许』」的邻域，**需要人工批准**；不共享的代价＝两份枚举会漂。按票面 AC#3 此句**交编排者/owner 裁**，本腿只交读数与先例（今后宿主读面按「只读、独立实现、路径决策走 `risk.PathResolver`」先例办），不把它当定案。

---

## 4.（附）AC#4 与票 181 的分工——顺带核毕

git **状态显示**＝票 181（`-done` 7/7，`git.go:144-171` 十一枚 json）；git **内容审查**＝票 189（`-done`，`Uncommitted`/`Diff` 在同一枚 `GitView` 的嵌套子段 `:178`/`:181`，无顶层新键）。两票都已结案，「互相归口到对方已勾格充数」今天无处可归；本程没做、也没建议这种归口。

---

## 5. 雷区清单（复核，沿用 r2 现数；一律单列待人拍板，不进只读交付）

现状合规锚现读：面板入向方法名册 `bridge.go:42-45` 四枚 `panel.*`＋`:66-67` `config.get/set`＋`:75` `config.models.list`；写动词点名尺 `grep -nE "\"(decide\|allow\|approve\|commit\|push\|revert\|stop\|remove\|apply\|toggle\|terminal)" internal/panel/bridge.go`＝**1 行且是 `:65` 注释逐字引禁令**；`lockedFieldFamilies`（`config_handlers.go:109`）逐字含 `plugins.` 族。

| # | 雷 | demo 控件现读 | Go 侧今天的门 |
|---|---|---|---|
| M1 | 附件移除 | `rightbar.js` `data-attach-del`＝**2** | 入向无 remove |
| M2 | diff 接受/应用 | `rb-review.js` `rb-apply-btn`＝**2** | 入向无 apply/accept |
| M3 | 插件启用开关 | `rb-plugins.js` `data-plg-toggle`＝**2** | 入向无 plugin 方法；`config.set` 锁定 `plugins.` 族 |
| M4 | 面板里发起执行 | `rb-terminal.js` `data-terminal-add`＝**2** | 地基都不存在（注册循环 5 支无 Shell） |
| M5 | 面板「允许」（潜在） | demo 无此控件：`rightbar.js` 允许/拒绝/批准/approve/allow 尺＝**2 行皆散文** | 入向无 decide/allow；`bridge.go:65` 注释引禁令 |

**现数 = 4 枚现存（M1–M4）＋1 枚潜在（M5）**，与 r2 同数同成员。本程没提「顺手把批准接上」。

---

## 6. 本件相对 r2 的订正（仅 1 处；其余全部复现）

| # | r2 断言 | 今天的尺 | 读数 | 判 |
|---|---|---|---|---|
| ① | r2 §2.1/§3 记票 92「6 勾」 | `grep -n "\[x\]\|\[ \]" .scratch/wisp/issues/92-…-done.md` | **5 勾／1 未勾**（未勾＝AC#5 变异三向；5 枚已勾＝AC#2/3/4/6/7）；该文件自 r2 锚以来 `git log --`＝**0 笔新提交**（最后两笔 `93bbb54d6`/`6d22dd6c7` 都在 r2 之前） | **r2 计数伪差**（多记 1 枚勾），非树移动；**不影响任何三档定性**（92 仍是 `-done`、K1-A/B 仍「已闭」）。谁引用 92 的框态请按 5/1 |
| ② | r2 全部其余读数（10 堆成员、三档、框态 14 张票、返图前提、雷区 4+1、`go list` 0/1、`PLAN.md` 行锚 `:3499-3512`/`:2540`/`:2575`） | 本件 §1–§5 同尺现跑 | **逐条复现** | 无过期 |

返图前提复核（承 r2 §6）：195 仍 `blocked`（票面 `:3`「排在票 33 片 B 交件之后」，`ls` 现量 `33-panel-host-c27.md` 无 `-done`）；407 仍 Blocked by 前端 face；187/195 在本题**一堆都不持有**的结论维持。认领尺 `grep -rlniE "内置终端\|文件树\|任务监控" .scratch/wisp/issues/*.md`＝**8 文件**（77/182/186/188/189/190/191/196），与 r2 同批，无第 9 枚。

---

## 7. 门禁与树卫生

- 派单明示「GATES: none needed (read-only analysis)」⇒ build/test/d22scan **未重跑**；承用依据＝射程内生产代码自 r2 锚**字节不动**（§0 diff 证据）＋ r2 门禁读数在册。
- `git status --porcelain -- internal/panel cmd/wisp`＝**2 枚 untracked**：`cmd/wisp/panel_resident_323_windows_test.go`、`cmd/wisp/panel_resident_assembly_323_test.go`——即派单点名的**票 323 在飞件**（`cmd/wisp/panel_resident_323*` 族）。本程**零接触**；本题任何尺都不以 `cmd/wisp` 测试件为对象，红/绿与本腿无涉。
- 本程 go 命令只有 `go list`（读 import 面）×2 次，无 build/test/写。
- `frontend/**` 零读零写（r2 的那条越界自报**未重犯**）；`docs/reports/pending-and-issues.md` 未碰。

---

## 8. 缺陷与未判项（具名）

1. **〔未复跑，转记〕** r2 §9-5 的「K1-C 分子/分母口径没量完」本程同样没量完（口径是 167 AC#2 的活，不是普查的活）；本程只把分母候选锚位刷新为 `budgets.go:79/:50/:54`。
2. **〔转记〕** `182-a1`/`182-c1` 对 `rightbar.js:90`（K9 在 demo 时间线里的混写位）的出处，本程未逐字复跑，沿用「转记」标记。
3. **〔不可判→已收窄〕** r2 §9-2 的 188 AC#2 口径问题维持原判（`Record` 尺 0 命中＋正控 1），**本程不改票**。
4. **〔计数伪差在册〕** §6-①：92 的框态以 5 勾/1 未勾为准；r2 件相应两处「6 勾」作废。本件不改 r2 文件（证据件只增不改），订正以本节为准。

---

## 9. 一句话交件

**r3＝确认跑：10 堆（7 有位＋3 有名无位）同数同成员；三档定性、归属票、雷区 4+1、`go list` 0/1、返图前提全部逐条复现 r2；唯一订正＝票 92 框态实为 5 勾/1 未勾（r2 多记 1 勾，不影响定性）。级联立票 0 枚（10 堆全有主；雷区按三程一致纪律继续单列待人拍板）。零产码、零翻勾、零台账、`frontend/**` 零读零写、既有票面零字节；写面只有本文件（新建）。**
