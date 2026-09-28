# 180-c1 ＋ 182-c1 — 面板那一维的 Go 侧现量普查（只读·零产码）

> 性质：**只读普查**。本件只交「界面上那一维，Go 里到底有没有真值、有没有读者、归属哪枚票」的
> 现量与代价。**不勾任何 AC 框、不动任何产码、不选落点、不提"顺手把批准接上"**
> （由面板侧来源的 L2「允许」是已定案禁止项，`AGENTS §1.2` ＋ 台账 `Q-49` 那一族）。
> 权威来源：`.scratch/wisp/dispatches/2026-09-28-145x-readonly-180-c1-plus-182-c1-panel-field-census.md`
> ＋ 票 180 ＋ 票 182。本件不新造规矩。

## 1. 起手锚＋写面闸门

- 起手时刻：`2026-09-28 14:40 +0800`
- 起手 HEAD：`e9ef94d0`（Mon Sep 28 14:35:19 2026 +0800）
  ⚠ 派单写的编排者锚点 HEAD＝`03c01d71`；起手实测 HEAD 已是 `e9ef94d0`（共享工作树，别人先走了）。
  本程所有行号一律**现跑**，不抄任何票面历史号。
- 起手写面闸门：`git status --porcelain -- internal/ cmd/` ⇒ **空输出**（干净）。
- 终态写面闸门：见 §10 末（待落）。

## 2. 票 180 AC#1 — 面板窗口尺寸归属三问

三把尺全部现跑（锚 `eba0b89b` 之后的工作树，本程 14:4x）。

**① 配置里那枚字段**（尺 `grep -n "width\|Width" internal/config/schema.go` ⇒ 命中 2 行）
- `internal/config/schema.go:527-528` 逐字：
  `// Width is the panel width in px.` ＋ `Width int \`toml:"width" default:"640"\``
- 所在结构体：`PanelSection`（`internal/config/schema.go:525-535`，五枚键全在场：
  `Enabled` `Width` `Height` `KeepAliveInSession` `Scale`）。

**② 生产码里谁读它**（尺 `grep -rn "\.Width" --include=*.go internal/ cmd/ | grep -v _test.go`）
- ⇒ **空输出（rc=1）**，命中 **0 枚**。票面 `039efb47` 的"读者＝零枚"到今天**仍然成立**。
- 顺手把同段另三枚也跑了：`grep -rn "\.Height\b\|\.Scale\b" --include=*.go internal/ cmd/ | grep -v _test.go`
  ⇒ **唯一命中 `internal/agent/budgets.go:94` 的 `b.Scale`**，那是 `CostSection` 的预算缩放，
  与 `PanelSection.Scale` 无关 ⇒ `Panel.Height`／`Panel.Scale` 生产读者同样 **0 枚**。
- `PanelSection` 整体在生产里**唯一**被碰到的地方是热重载的**通用拷贝**：
  `internal/config/manager.go:211` 逐字 `{"panel", &cur.Panel, &fresh.Panel, func() { cur.Panel = fresh.Panel }},`
  ⇒ 它按指针搬整枚 struct，**不读 `Width` 的值、也不产生任何效果**。

**③ 窗口／WebView 的实际尺寸今天从哪来**
- 尺 `grep -rn "NewWindow\|SetBounds\|Rect{\|Width:" --include=*.go internal/panel/ | grep -v _test` ⇒ **空输出**。
- 尺 `grep -rn "WebView2\|CreateWindow\|bounds\|Rect" --include=*.go internal/panel/ internal/winsec/ cmd/ | grep -v _test.go`
  ⇒ 命中只有两类：(a) **注释**说这棵树里根本没有宿主；(b) `cmd/wisp/notify_windows.go:149` 的
  `pCreateWindowExW.Call(...)`——同一函数 `:153` 的错误串逐字 `CreateWindowExW(message-only): %v`，
  即那是**通知用的 message-only 隐藏窗口**，不是面板。
- 注释现量（逐字，三处）：
  - `internal/panel/doc.go:16`：`// DEFERRED(host/bridge): implemented by ticket 33 (host), ticket 35 (bridge).`
  - `internal/panel/pump.go:16`：`// tree today - no WebView2 host (ticket 33 is unclaimed), no postMessage writer,`
  - `cmd/wisp/run.go:438`：`// have is a page to go to: this tree carries no WebView2 host (tickets 33/35),`
- ⇒ **今天没有任何人决定面板尺寸，因为面板窗口本身不存在**（票 33 的 host、票 35 的 bridge 都还是
  `ready-for-agent`；`ls .scratch/wisp/issues/` 里 `33-panel-host-c27.md`／`35-panel-bridge-c17.md` 均**无 `-done` 后缀**）。

**②→③ 之间缺的不是"一跳"，是整条链**：`cfg.Panel.Width` → 宿主 → `SetWindowPos/尺寸`。
中间那一环（WebView2 host）从未落地，所以 `Width` 与 `Height`/`Scale` 一样是**悬空字段**。

**三档结论：＝「规格要求生效」那一档**（⇒ 按票 180 AC#2，正解是把那一跳接上；但**接哪一跳归编排者选落点，本程不选**）。
规格原句（不转述，逐字）：
- `docs/PLAN.md:2726`：`**三档生效级别**：\`hot\`（立即生效）· \`reload\`（需重载子系统，如换 ASR 模型）· \`restart\`（需重启进程）。`
- `docs/PLAN.md:2743`（D36 配置树那一行）：
  `| \`[panel]\` | \`enabled\` \`width\` \`height\` \`keep_alive_in_session\`(true) \`scale\` | \`hot\` |`
- `docs/specs/SPEC-03-config-secrets-envs.md:39`：
  `| \`[panel]\` | \`enabled(bool)=true\` \`width(int)=640\` \`height(int)\` \`keep_alive_in_session(bool)=true\` \`scale(float)\` | hot |`

⇒ `width` 被明确列进 `hot`，而 `hot` 的定义逐字就是「立即生效」⇒ **规格要求它生效，本票属"补实现"支，不是"改规格文字（＝人工批准）"支**。
反向自查（规格有没有另一处只要求"显示"）：尺 `grep -rn "width\|Width" docs/specs/SPEC-03*.md docs/specs/SPEC-08*.md`
⇒ SPEC-08 里唯一命中是 `SPEC-08-ui-ball-panel.md:218` 的 SVG `stroke-width`，**与面板窗口无关**；
`SPEC-08` 全文**没有一处**给面板窗口定尺寸。

## 3. 票 180 AC#4 — 同族「带 `default` 但生产零读者」枚数与逐枚名

尺（现跑，逐枚）：对 `internal/config/schema.go` 里每一枚带 `default:` 的字段名跑
`grep -rn "\.<字段名>\b" --include=*.go internal/ cmd/ | grep -v _test.go`，取命中数。
原始读数落 `.scratch/wisp/probes/180/c1/same-family-census.txt`（61 行）与
`probes/180/c1/field-to-struct-map.txt`（69 行，字段→struct 映射）。

- 基数：`schema.go` 里带 `default:` 的**行＝69**，去重后的**字段名枚数＝61**；
  非测试 Go 文件总数＝231（`internal/` ＋ `cmd/`）。
- **零生产读者（total=0）＝18 枚**，逐枚名（限定名，映射表现跑）：

| # | 限定名 | 所在 struct |
|---|---|---|
| 1 | `PanelSection.Width` | `PanelSection`（票 180 本尊） |
| 2 | `PanelSection.KeepAliveInSession` | `PanelSection` |
| 3 | `BallSection.ClickThrough` | `BallSection` |
| 4 | `BallSection.HideOnFullscreen` | `BallSection` |
| 5 | `SessionSection.WarmTimeoutSec` | `SessionSection` |
| 6 | `SessionSection.SettlingSec` | `SessionSection` |
| 7 | `SessionSection.ConversationIdleSec` | `SessionSection` |
| 8 | `AudioSection.InputDevice` | `AudioSection` |
| 9 | `AudioSection.SampleRate` | `AudioSection` |
| 10 | `PrivacySection.DiagnosticsOptIn` | `PrivacySection` |
| 11 | `PrivacySection.RetentionDays` | `PrivacySection` |
| 12 | `MemorySection.L1Enabled` | `MemorySection` |
| 13 | `MemorySection.L1Max` | `MemorySection` |
| 14 | `MemorySection.L3RetentionDays` | `MemorySection` |
| 15 | `AECConfig.EchoRef` | `AECConfig`（`voice.aec.echo_ref`） |
| 16 | `RollConfig.SizeMB` | `RollConfig`（`observe.roll.*`） |
| 17 | `RollConfig.Days` | `RollConfig` |
| 18 | `ObserveSection.SLOSampleIntervalSec` | `ObserveSection` |

- ⚠ **18 是下界不是全量**：粗尺 `\.<名>` 对**同名异处**的字段会**多报**读者
  （`Enabled`/`Mode`/`Size`/`Level`/`Max`/`Provider` 这类泛名，命中很可能来自别的 struct），
  所以"非 0"不等于"这一枚有读者"；反过来"0"是**铁的零读者**。⇒ 真实枚数 **≥18**。
  票 83 那张表用 I-A~I-G 七把仪器处理的正是这个盲区（见下）。
- **逐枚归口＝不用另立名单：这 18 枚全部已在「票 83（done）」的 AC#1 全量键表里登记过**。
  尺 `grep -rliE "<键名>" .scratch/wisp/issues/*.md` ⇒ 命中 `83-config-keys-that-lie-must-fail-loudly-done.md`，
  逐字（`issues/83-...md:155`）：
  `| \`panel.{enabled,width,height,keep_alive_in_session,scale}\` | \`schema.go:517-525\` | 无（I-A） | 📋 票 33/34/36 |`
  其余各枚同样在表内有行：`:131` ball 两枚 → 📋 票 62/65/68 ＋ 票 39；`:133` session 三枚 → 📋 票 28；
  `:134` `voice.*`（含 aec）→ 📋 票 15/26/27/41/59/60/61；`:135` audio 两枚 → 📋 票 13 尾巴／票 26；
  `:152` privacy 两枚 → 📋 票 45／票 08；`:154` memory 三枚 → 📋 票 29；`:161` observe/roll 三枚 → 📋 票 08/42/66。
  ⚠ 该表引的 `schema.go:517-525` 等行号**已被注释级改动推漂**（今天 `PanelSection` 在 `:525-535`），
  本程一律用现跑号，不沿用票 83 的历史号。
- ⇒ **对票 180 现量第 3 条的复核结论**：`internal/config/unwired_test.go` 的"名册"**射程只有四枚 locked section**——
  `unwired_test.go:311 TestEveryLockedSectionKeyIsAccountedFor` 里 `collect(...)` 只喂
  `risk` / `fs` / `net` / `plugins`（`grep -n "collect(" internal/config/unwired_test.go` 现量），
  守卫表 `lockedKeyDisposition` 在 `internal/config/unwired.go:119`；`:217` 那行 `width = 800` 属
  `TestUnwiredGuardIsScopedNotABigStick` 的"诚实配置不该报错"样本，**不是登记**。
  ⇒ `[panel] width` 不在名册射程内是**设计如此**，本票第二格因此**既不是"补登记"、也不是"新发现"**：
  发现与归口都已在**票 83 那张表**里；真正缺的是**把效果接上**那一格（AC#2，落点归编排者）。

## 4. 票 182 AC#1 — 「任务监控」那一栏的堆数与逐堆名

**堆数（现量）＝7 堆。** 每一堆都必须写"我在谁那儿看到的这一堆"；
⚠ **owner 给的截图是外部产品（qoder）的界面，不是本仓规格**——本表只按仓内文本现量，
凡只出现在票 182 票面"目测清单"里的那一堆，都明写成"目测、无规格"。

| # | 堆名 | 我在谁那儿看到的这一堆（现跑出处） |
|---|---|---|
| S1 | 子代理／后台任务状态 | 票 182 AC#1 目测清单；"子代理"在产品侧**仓内零命中**（尺见 §5），issue 池里的"子代理"全指开发分工不是功能 |
| S2 | 环境信息（本地路径／分支／未提交枚数／提交或推送） | 票 182 AC#1 目测清单 ＋ 票 181 票面（分支／工作树）＋ 票 186（composer 那一排的 git，**不同地界**） |
| S3 | 代码审查 diff | 票 77:364 逐字"`5e23d99`（右栏三标签定案：审查/终端/浏览器）"、票 77:384 逐字"审查/终端/浏览器三面板、任务监控弹窗"；＋票 182 目测 |
| S4 | 工作区文件树 | **只在票 182 票面的目测清单里**；尺 `grep -rln "文件树" .scratch/wisp/issues/*.md` ⇒ 票 77／票 186 里**没有**这三个字（本程 14:5x 现跑，命中集合见 `probes/182/c1/stack-rulers-precise.txt`） |
| S5 | 内置终端 | 票 77:364／:384（右栏标签之一）＋票 182 目测＋D34 表 `PLAN.md:2569` 的 `shell.session` 行 |
| S6 | 浏览器 | 票 77:364／:384（右栏第三标签）；票 182 目测清单**没有它** ⇒ 本程按仓内文本补进堆数 |
| S7 | 成本与 token | 票 182 目测清单 ＋ **票 145 那张十四行表的第 12 行**（`docs/evidence/s1/145-snapshot-field-census-r1.md:182` 逐字 `| 12 | :3486 | 成本 | \`Cost{In,Out,Cached,Micros,Currency}\` | ✅ 五枚全有源 |`；该行对应的 `PLAN.md` 真身今天在 `:3492`，票 145 引的号已漂 6 行） |

堆数口径：S1–S7 是我能从**仓内文本**逐枚点名的堆；owner 目测清单里那六项**全部落进**了 S1–S7，
没有被丢掉或被合并。**没有一列"其它"**——那一栏里若还有别的堆，需要下一位以仓内文本为准补点名。

## 5. 票 182 AC#2 — 逐堆三档定性（每档一条现跑尺）

尺与命中数原始读数：`.scratch/wisp/probes/182/c1/stack-rulers.txt` 与 `stack-rulers-precise.txt`。
正控（证明尺本身能打数）：同一批调用里 `任务/会话状态机` 那把尺命中 **34**、
`成本 token 用量` **12**、`fs list / tree` **11** ⇒ **"0 命中"这几发不是尺坏了**。

| 堆 | 三档 | 现跑尺 ⇒ 命中 | 取到真值要从哪枚函数走 |
|---|---|---|---|
| S1 子代理 | **规格真空**（别造） | `grep -rnE "(Subagent\|sub_agent\|SpawnAgent\|DelegateTask\|fanout)" --include=*.go internal/ cmd/` ⇒ **0** | 无源可走：产品里今天**没有**子代理这一层 |
| S1 后台任务 | **有源**（缺载体） | `grep -n "type TaskRoster" -A / Roster` ⇒ `internal/tools/task.go:110` 在场，生产写入者 `task_backfill.go:115` | `TaskRoster.Record/Look/Count`（`internal/tools/task.go:123/139/152`）；装配点 `cmd/wisp/run.go:364`＋`:639`。⚠ 它存的是 `TaskOutput`（输出续读用），**没有状态／耗时／名字** ⇒ 借它显示"任务在跑什么"要新加维度 |
| S2 本地路径 | **有源**（已有载体） | `grep -rn "type WorkspaceView" -A 12 internal/panel/` ⇒ `composer.go:174-186`（`Set/Spelling/Canonical/Reparse/Rewritten`） | `panel.NewComposerState`（`composer.go:210`）→ 已在 `Snapshot.Composer.Workspace`；落地者＝票 92（`-done`） |
| S2 分支／工作树 | **无源但已有人要求**（票 181，非冻结契约） | `grep -rnE "\b(HEAD\|branch\|worktree)\b" --include=*.go internal/panel/ internal/config/ cmd/` ⇒ 6 命中**全是注释里的 control-flow "branch"**，真 git 读面 **0** | 无：今天没有任何函数读 `.git/HEAD`。唯一 git 命中来自 `internal/risk/blacklist.go:214-217,329`（`.git-credentials`／`.git/config` 的**拉黑路径分类**，是判级不是读面） |
| S2 未提交枚数 | **规格真空**（票 182 AC#4 把它划给"本票那一堆"，但本票只普查不实现） | 同 S3 尺 ⇒ **0** | 无 |
| S3 审查 diff | **无源＋规格真空** | 尺（票 182 现量第 4 条同形）`grep -rn "git diff\|diff --numstat\|numstat\|unified diff" --include=*.go internal/ cmd/` ⇒ **0** | 无 |
| S4 文件树 | **模型侧有源、宿主侧读面无源；面板那一维规格真空** | `grep -rn "FSListDecl\|fs\.list" --include=*.go internal/ cmd/` ⇒ 11 命中，全在 `internal/tools/fs.go`（`:190 Name() "fs.list"`、`:307 FSListDecl`）与注释 | 真值走 `fs.list` 的执行体（`internal/tools/fs.go`），但它带 **D34 风险分级**（`PLAN.md:2534` 逐字 `| \`fs.list\` | 列目录 | L0 / L2（越界） | \`fs.read\` | S3 | — |`），且是**模型调用面**不是宿主读面 ⇒ 见 §6 那句"要不要共享实现" |
| S5 内置终端 | **工具那一维＝无源但规格要求（D34 表内行）**；**面板载体＝规格真空** | `grep -rnE "\b(PTY\|pty\|ShellSession\|shell\.session)\b\|cmd\.exe" --include=*.go internal/ cmd/` ⇒ **1 命中且是注释**（`internal/risk/rules_shell.go:10` 讲 cmd.exe 转义），PTY 零命中 | 无：`shell.session` 在 D34（`PLAN.md:2569` 逐字 `| **\`shell.session\`** | **在一条活着的会话里跑下一条命令（工作目录与环境延续）** | **L2** | \`shell\` | S3 | 票 163（20…`）但**生产注册表里没有它**。
  **生产在册名册现跑（订正版）**：注册只有两个循环——`cmd/wisp/run.go:345` 的
  `for _, e := range tools.BuiltinFSEntries(tools.FSDeps{`（`internal/tools/fs.go:325-331` 逐字两枚
  `fsRead`/`fsList` ＋ `return append(out, BuiltinFSWriteEntries(d)...)`）与
  `cmd/wisp/run.go:365` 的 `for _, e := range tools.BuiltinTaskEntries(tools.TaskDeps{Roster: rt.tasks})`
  （`internal/tools/task.go:295-297` 逐字一枚 `task.output`）。
  ⇒ 本程把 **`fs.read`／`fs.list`／`task.output` 三枚核到名册级**；
  `fs.edit`/`fs.delete`/`fs.move`/`fs.trash`/`fs.write` 那几枚在 `BuiltinFSWriteEntries` 那一支里，
  **本程没逐枚展开**（写它们需要另一次读数）。
  ⚠ 早前本程一次粗尺把 `probe.count` 当成在册名——**订正：它是 `internal/tools/bridge_test.go:233` 的测试替身，不在生产**。
  `shell.exec`／`shell.session` **不在任何一支里** ⇒ 与票 163 AC#1b 的读数一致。 |
| S6 浏览器 | **无源**；⛔ **实现路径＝未定义即停项，本程不判** | `grep -rnE "browser\|Browser" --include=*.go internal/ cmd/` ⇒ 6 命中全在 `internal/risk/blacklist.go:227-239`（浏览器凭据库拉黑）与 `cmd/wisp/panel_assets.go:8` 注释 | 无。**`web.search` 实现路径是 `AGENTS §2` 逐字列出的"未定义即停"项（S3 前）** ⇒ 碰到就停手，不选路 |
| S7 成本／token | **有源**（缺载体） | `grep -rnE "(CostOf\|TrackCost\|total_tokens\|TotalTokens\|microUSD\|MicroUSD)" --include=*.go internal/ cmd/` ⇒ **12** 命中 | `internal/agent/cost.go:38 AddUsage(p config.Price, u llm.Usage) int64` → `guard.go:155 AddCost` → `loop.go:976 publishUsage`；实时 token 另有 `llm.StreamEvent.Usage`（票 145 表 §2.13 已点名）。⚠ 票 145 那行还挂着**单位口径未定案**（`config.Price` 是 micro-USD，`task_log.currency` 默认 CNY） |

## 6. 票 182 AC#3 — 逐堆归属指认表（含「没票认领」那一列）

⚠ 规矩：指认前**逐枚现跑**过那张票的 AC 标题（`grep -n "AC#" <票文件>`），命中的原文列在下面；
近亲票（186 的 git 切换、187 的模型档位）**不拿来填空**。

| 堆 | 归属 | 现跑凭据（那张票的 AC 逐字／或其反证） |
|---|---|---|
| S1 后台任务状态 | **没票认领 ⇒ 要立票** | 票 145 十四行表**逐行点名**只有：思考中／SSE 流式／推理过程／工具调用／工具调用（展开）／审批等待／L2 确认卡／L1 阻止窗口／L2 原生降级卡／错误／Stuck／成本／已取消／注入检出——**没有"后台任务"这一行**（尺：`sed -n '168,192p' docs/evidence/s1/145-snapshot-field-census-r1.md`） |
| S1 子代理 | **没票认领 ⇒ 要立票，且先得有规格**（规格真空） | issue 池全文尺 `grep -rlniE "任务监控\|文件树\|内置终端" .scratch/wisp/issues/*.md` ⇒ 只 3 枚（182／186／77），"子代理"在产品侧零命中 |
| S2 本地路径 | **已有票落地**（票 92，已 `-done`） | `92-panel-composer-mode-attachments-workspace-done.md` AC#3 逐字"工作区选择：一次真实切换后，**(i)** 后续操作的风险判定用的是新工作区" |
| S2 分支／工作树（只读显示） | **票 181** | 票 181 AC#1 逐字"① 当前分支（`.git/HEAD` 的 `ref:` 行…）；② 本地分支列表＋远端分支名…③ 已有工作树（`.git/worktrees/` 目录枚举）" ⇒ 这一堆**真写着这件事** |
| S2 未提交枚数／提交／推送 | **没票认领 ⇒ 要立票**（＋动作那一半是雷区，见 §7） | 票 181 AC#4 只管"非 git 目录那一形"的显示；票 186:48 逐字"不动输入框宽度（＝票 180）、**不裁任务监控那一栏缺哪些堆（＝票 182）**" ⇒ 186 明确不认领 |
| S3 审查 diff | **没票认领 ⇒ 要立票** | 票 181 AC#3 逐字"交付里**不许出现**任何模型可调用的 `git.*` 工具…`internal/panel/bridge.go:42-45` 那四枚 `panel.*` **一枚不许加**" ⇒ 181 只走只读快照通道，不含 diff |
| S4 文件树 | **没票认领 ⇒ 要立票** | 无一张票的 AC 写着"面板显示工作区树"；`fs.list` 的票面归口是 S3 工具切片（D34 行 `PLAN.md:2534`），不是面板载体 |
| S5 内置终端（工具地基） | **票 163（只登记依赖，本程不重开）** | 票 163 AC#1b 逐字"`shell.exec` 与 `shell.session` 两枚名字在生产注册表里都不存在…差集非空就不许进实现格"（与 §5 现跑在册名册一致） |
| S5 内置终端（面板载体） | **没票认领 ⇒ 要立票** | 票 163 的 AC 全是工具/会话本体，无"面板"字样 |
| S6 浏览器 | **没票认领（Go 侧）**；前端标签在票 77 | 票 77:364／:384 是**组件清单不是 Go 数据源**（票 182 现量第 1 条同判）；⚠ 实现路径撞 `AGENTS §2` 未定义即停项 |
| S7 成本／token | **票 145 第 12 行**（落地格 AC#2 未勾）＋ 配额那一层票 44 | 票 145 表行 12 逐字见 §4 表末；票 83 AC#1 表 `:156` 逐字 `| \`cost.{daily_budget,monthly_budget,alert_threshold,over_budget}\` | … | 📋 票 44（C23） |` |

**AC#3 特别要答的那一句（宿主侧读面 vs 模型侧工具要不要共享实现）——本程只报现状、交编排者裁**：
- 现状：仓里已经有**两次**把边界写成"只读显示不走模型工具"的定案文本——
  票 181 AC#3（"只读显示走**已有的快照推送通道**，不需要新请求方法"）与
  票 186 AC#5（"不许把切换能力塞进 `internal/tools/` 的执行面"）。
- 风险点：`fs.list`／`task.output` 都是**带 D34 风险分级的模型调用面**（`fs.list` 是 `L0 / L2（越界）`）。
  面板若直接复用同一执行体，就把这两级**渗进面板通道**（反向是把面板变成执行面）；
  若另写一份宿主侧读面，则同一目录列表**两份实现**会漂（票 146 那族"显示与真用分账会静默漂"同形）。
  ⇒ **两条都有代价，本程不选，交编排者裁**。

## 7. 票 182 AC#5 — 雷区清单（逐堆：有无面板发起的批准／写动作）

禁令逐字（`AGENTS §1.2` 抄自 SPEC-06/票 141 交付版）：**"由面板侧来源的 L2『允许』"** 是禁止项；
已定案口径逐字（`docs/PLAN.md:2027`）：
`① **L2 的「允许」决策不接受来自面板的调用**（\`approval.decide\` 服务端直接拒绝 panel 来源的 allow）`
`② 面板只提供**「拒绝」与「查看完整参数」**…③ 批准只能由原生侧产生：**悬浮球点击 / 原生确认卡按钮 / 全局快捷键**`

| 堆 | 有无面板发起的批准／写动作 | 现状凭据（现跑） |
|---|---|---|
| S1 后台任务状态 | **无**（只读） | `Snapshot` 四枚里没有任务字段（`composer.go:57-62`） |
| S2 环境信息 | ⚠ **有（在"提交或推送"那一项上）** | owner 目测清单里那一支就是写动作；`git` 子进程今天全仓零次（§5 尺）；票 181 票面标题逐字"…the **two mutating actions** would need C17 whitelist methods" ⇒ **待人拍板**，不许并进只读那一堆 |
| S3 审查 diff | 只读那一半**无**；⚠ "接受／回滚改动"那一半＝写动作（`fs.write`／`fs.edit` 级，L2） | 今天 diff 零源 ⇒ 两半都不存在；登记为**立票时的边界** |
| S4 文件树 | 只读**无**；⚠ 树上带删除／移动＝写动作（`fs.delete`／`fs.trash`／`fs.move` 在 `internal/tools/` 有 `Name()` 声明，D34 判 L2；⚠ **是否真在生产名册本程没逐枚展开**，见 §5 订正行） | 生产在册名册只核到 `fs.read`／`fs.list`／`task.output` 三枚（§5） |
| S5 内置终端 | ⚠ **最大雷**：面板里敲命令＝执行面进面板通道；批准来源问题直接撞 `PLAN.md:2027` ①③ | 票 163 AC#4 逐字"常驻会话**不许**变成'授权一次就一直放行'…不许动 `internal/agent/approval/**`" |
| S6 浏览器 | 只读浏览无；⚠ 若带"在页面里执行/登录"＝撞 `AGENTS §1.2` 的"只靠 CSP＋净化的方案已被否决"（`PLAN.md:2027` ⑤） | 零源 |
| S7 成本／token | **无**（只读） | `publishUsage` 是出站播报，不是入口 |
| （参照格）审批卡 | **现状合规**：面板四枚方法里没有 `decide`／`allow` | `internal/panel/bridge.go:42-45` 逐字四枚＝`panel.mode.request`／`panel.workspace.request`／`panel.attachment.add`／`panel.message.send` |

⇒ **本程不提"顺手把批准接上"**；上表 4 处 ⚠ 一律**单列为待人拍板项**（见 §12）。

## 8. 面板要的字段清单（交编排者转前端会话；本程不写 `frontend/**`、不裁 TS interface）

枚数＝**8 枚**（A 组 1 ＋ B 组 7）。来源逐枚标注"已存在／要新写／规格真空"。

| # | 字段（Go 侧候选名） | 属哪堆 | 来源档位 |
|---|---|---|---|
| A1 | `Panel.Width`（＋同段 `Height`/`Scale`/`KeepAliveInSession`） | 票 180 | **字段已在、值无读者**；`Width`/`KeepAliveInSession` 带 `default`，`Height`/`Scale` 连 `default` 都没有（`schema.go:525-535`）⇒ 生效链整条缺（宿主不存在） |
| B1 | `Run.Phase`／后台任务状态 | S1 | **无源**（子代理）／**有源缺载体**（`TaskRoster`，但没有状态维） |
| B2 | `Env.Branch`／`Env.Worktrees` | S2 | **无源**，票 181 已要求（AC#1 三支） |
| B3 | `Env.UncommittedCount` | S2 | **无源＋规格真空** |
| B4 | `Review.Diff` | S3 | **无源＋规格真空** |
| B5 | `Tree.Entries` | S4 | 模型侧有源（`fs.list`）／宿主侧无源＋规格真空 |
| B6 | `Terminal.Session` | S5 | **无源**，工具地基＝票 163（在册零实现） |
| B7 | `Cost{In,Out,Cached,Micros,Currency}` | S7 | **五枚全有源**（票 145 行 12），缺载体＋单位口径未定案 |

## 9. 本程没测什么（逐名）

1. **没跑过任何真窗口**：面板宿主不存在（票 33/35 未开工），所以"改 `config.toml` 界面动不动"这一发
   今天**根本没有可观测面**——本程只交静态现量，没做真机验证，也做不了。
2. **没测 `frontend/**`／`design/**`**：按派单禁令不读不写不引不转述 ⇒ §8 那份字段清单是**Go 侧**清单，
   不是"前端应当长这样"的 TS 契约（票 145 AC#2 那一格明确把 TS 对齐留给前端会话）。
3. **`unwired_test.go` 只读了语义与 `collect(...)` 的四枚 locked section**，没单独跑它的用例
   （跑数在 §10 的门禁里，不是本程的判据）。
4. **粗尺 `\.<字段名>` 的同名歧义没逐枚消**：§3 因此只报**下界 18 枚**；
   对"非 0"的 43 枚没做"这一枚到底有没有读者"的二次判定（那是票 83 用 I-A~I-G 七把仪器做过的事，
   本程复用其结论、不重造）。
5. **没裁数值**：`width` 该是多少（640?）＝owner 的一句话，本程不碰、也不替 owner 答。
6. **没选落点**：`Width` 生效链要接在哪一枚（票 33 的 host？新票？票 145 的快照？）＝编排者的活。
7. **S6 浏览器那一堆的实现路径没判**：撞 `AGENTS §2` 的"未定义即停"项（`web.search` 实现路径，S3 前）。
8. **S1 子代理这一维在规格里没有对应物**：本程只报"零命中"，没有替它造规格。
9. **§4 那句"仓内文本之外还有没有别的堆"没穷尽**：堆数＝7 是**能从仓内文本逐枚点名**的量，
   不是那一栏的像素级清点（像素级归前端会话）。

## 10. 门禁终态

⚠ 按 `A363` 薄规矩：**只读程取的这些数不充当任何 AC 的结案凭据**。

**0. 重跑前先查写数落点（`A367` 那一味）**
尺 `grep -n "WriteFile\|OpenFile\|> \"\|>> \"\|tee" scripts/d22scan.sh .scratch/wisp/probes/154/gate-clauses.sh` ⇒ **空**；
再尺 `grep -nE "^\s*[a-z_]*\s*>|>>|cp |mv |tee|git checkout" .scratch/wisp/probes/154/gate-clauses.sh` ⇒ 只有 4 行
**进程内 `printf`**（`:82/:180/:192/:193`，往 `grep` 喂字符串）。
⇒ 两枚台件**都不落盘**，重跑不会洗掉别家票的原始读数；本程自己的原始读数一律写到
`probes/180/c1/**` 与 `probes/182/c1/**`。

**1. `sh scripts/d22scan.sh` ⇒ rc=0（＝基线）**
逐字读数：`d22scan: examined 231 production Go files under internal/ and cmd/`、
`d22scan: scope ban #8 internal/  examined 433 Go files, comments and _test.go included`（**433＝派单给的基线**）、
`ban #8 cmd/ examined 45`、末行 `clean - no D22 ban violations`。

**2. `bash .scratch/wisp/probes/154/gate-clauses.sh` ⇒ 比名册不比退码**
`腿数＝14 声明与实测不符＝1`；`名册=14 声明=14 记账=14 缺腿=0 空头声明=0`；`基线过期枚数＝0`；`聚合退码＝1`。
**红腿名册＝只有 `G6neg`**，逐字 `# BAD  腿=G6neg 声明=ring 基线=1枚 实测=3枚 因=新增未成对（票 171 AC#2：实测 > 基线）`
⇒ **与派单给的在册名册（`G6neg`＝票 178）相符，差集为空**。
⚠ 那枚 `实测=3 > 基线=1` 的增量**不是本程造的**（本程零产码，`git status --porcelain -- internal/ cmd/` 见 §1 与文末闸门）。

**3. `go test -count=1 ./internal/config/` ⇒ `ok github.com/CarlosShao/wisp/internal/config 0.790s`（rc=0，0 红）**

**4. `go test -count=1 ./internal/panel/` ⇒ FAIL，红 **2** 枚（逐名）**
- `--- FAIL: TestPanelColourLiteralsLiveOnlyInTheGeneratedTheme`
  （`frontend_hygiene_test.go:219`：`57 files reachable from main.tsx carry colour literals only in the generated theme`）
- `--- FAIL: TestC21DesignTokensFourWayAgree`
  （`tokens_fourway_test.go:441`：`read design/assets/tokens.css: open ...\design\assets\tokens.css: The system cannot find the path specified.`）
⇒ 两枚红的射程都在**本程禁令之内没碰的东西**上（`frontend/**` 与本机缺失的 `design/assets/tokens.css`），
与票 77 AC#2 在案的同名红（"Go 侧 `tokens.go` 与 `c21-native-tokens.md` 未同步"）同族；
本程**没动 `internal/**` 一字**，所以这两枚不是本程造的，也**不是本程能修的**。

**5. 没跑 `./cmd/wisp/`**：本机测不到东西＝票 98，按派单不算进本程的红。
**6. 名册两向 `comm` 差集**：本程零产码 ⇒ 没有新增/删除的在册名可差，**未跑**（不是漏跑，是没有对象可差）。

**终态写面闸门（两条，起手＋终态）**
- 起手 14:40：`git status --porcelain -- internal/ cmd/` ⇒ **空输出**。
- 终态：见文末「终态补记」，取在最后一枚 commit 之后。

## 终态补记

（取在本程最后一枚 commit **之后**现跑，14:5x）

- 本程 commit 链（5 枚，全部显式 pathspec、只 commit 不 push）：
  `eba0b89b` 骨架 → `6febca94` 票 180 §2/§3 → `ad759e94` 票 182 §4–§8 → `92a4f2d8` 门禁＋两枚票面 log＋台件＋名册订正 → 本件最后一枚。
  ⚠ 共享工作树：起手实测 HEAD `e9ef94d0`（不是派单写的 `03c01d71`），中途别人落了 `378c8b34`（票 174-r2）等；
  本程所有行号一律现跑，未沿用任何票面历史号。
- **终态写面闸门**：`git status --porcelain -- internal/ cmd/` ⇒ **空输出**（第二枚，与起手那枚成对，两条都空）。
- 终态门禁复测（最后一枚 commit 之后）：
  - `sh scripts/d22scan.sh` ⇒ 末行逐字 `d22scan: clean - no D22 ban violations`，
    `d22scan: scope ban #8 internal/ examined 433 Go files, comments and _test.go included`（**＝基线 433**），rc=0；
  - `go test -count=1 ./internal/config/` ⇒ `ok ... 1.047s`；
  - `./internal/panel/` 那 2 枚红见 §10 第 4 条（红因在 `frontend/**` 与本机缺失的 `design/assets/tokens.css` 上，本程没碰）。
- **零删除自证**：全程没有 `rm`/`del`/`git clean`/`git restore`/`git checkout .`/`git stash`/`git reset`/
  `git rebase`/`--amend`/worktree/switch；`probes/161/r6/flip-declaration.sh` 没跑；
  `probes/154/gate-clauses.sh` 只执行、未改一字。台件 5 份原始读数全部入库。
- **本程改过的文件全集**（四枚，均在派单 §1 允许的写面内）：
  `docs/evidence/s1/180-182-panel-fields-census-c1.md`、
  `.scratch/wisp/issues/180-...md`（只追加一段 Progress log，未改原句、未勾框）、
  `.scratch/wisp/issues/182-...md`（同上）、
  `.scratch/wisp/probes/{180,182}/c1/**`（新建 5 份读数）。
- **工具调用终值＝40 枚**（硬顶 45）；新探索停在第 30 枚，其后 10 枚全部用于写件、门禁与 commit。
  ⚠ 未判"必须改产码"⇒ 没有触发停手上报的门；两票的每一问都能在只读面上答完。

## 11. 被拒调用＋零删除自证＋工具调用终值

- **被拒调用：0 枚**（全程无一次权限被拒；无重试、无绕过）。
- **删除命令：一枚没跑**。全程没有 `rm`／`del`／`Remove-Item`／`git clean`／`git restore`／`git checkout .`／
  `git stash`／`git reset`／`git rebase`／`--amend`／worktree／switch；
  临时件只建不删：`probes/180/c1/{same-family-census.txt, field-to-struct-map.txt, issue-pool-claims.txt}`、
  `probes/182/c1/{stack-rulers.txt, stack-rulers-precise.txt}` 全部留在盘上并入库。
- **没跑 `probes/161/r6/flip-declaration.sh`**（禁令）；**没动 `probes/154/gate-clauses.sh` 那把尺本身**（禁令，
  只是执行它）。
- 重跑既有台件前先查了写数落点（`A367` 那一味）：现跑结果见 §10 第 1 行。
- **工具调用终值：见文末「终态补记」**（最后一枚 commit 之后现量）。

## 12. next＝落地腿还缺什么、哪几枚要人先批准

**A. 票 180（`Panel.Width` 生效）**——本程判＝**「规格要求生效」档**（凭据 §2 逐字三行），所以缺的是腿不是话：
1. **面板宿主本身**（`.scratch/wisp/issues/33-panel-host-c27.md` 无 `-done`，票 35 同）
   ⇒ 没有窗口就没有"生效"可言，**这是硬前置**；
2. 宿主落地后才有 `cfg.Panel.Width → 窗口尺寸` 那一跳，并要一发
   **能区分"接上了"与"又抄了一遍默认值"**的判据（票 180 AC#2 的口径）；
3. ⚠ `Height`／`Scale` **连 `default` 都没有**（`schema.go:530`／`:534`），比 `Width` 更空——
   要接就接一排；只接 `Width` 会留下三枚新装饰品（正是票 180 AC#3 反向判据防的那一形）。
⇒ **要人先批准：无**（`hot` 档已写在 PLAN/SPEC-03，接效果不动契约）。
但**落点选择**（并进票 33 还是立一枚"接线票"）归编排者，本程不选。

**B. 票 182（任务监控栏）**——「没票认领」那一列共 **7 处**需要立票或被点名：
S1 后台任务状态、S1 子代理（先得有规格）、S2 未提交枚数、S3 审查 diff、S4 文件树、
S5 面板终端载体、S6 浏览器（Go 侧）。
可并格的两处：S7 成本/token → **票 145 AC#2**（那格今天未勾）；S2 分支／工作树 → **票 181**。
S5 工具地基只**登记依赖**到票 163，本程不重开它。
⇒ **要人先批准（本程一律不选，只点名）**：
1. **雷区 4 处**（§7）：面板里的 git **提交／推送**、diff 的**接受／回滚**、文件树上的**删除／移动**、
   **面板内置终端**——每一处都连着"由面板侧来源的 L2『允许』"禁令（`PLAN.md:2027` ①②③逐字见 §7），
   **必须逐枚人拍板，不许并进只读那一堆**；
2. **宿主侧读面 vs 模型侧工具要不要共享实现**（§6 末，两条都有代价）；
3. **`web.search`／浏览器那一堆的实现路径**＝`AGENTS §2` 在册的未定义即停项；
4. **C17 白名单**：若面板终端或 git 动作要做成请求方法，`panel.*` 就要从四枚往上加
   ——票 181 AC#3 现文本是"**一枚不许加**"，动它＝**契约面＝人工批准**（`AGENTS §1.1`）。
